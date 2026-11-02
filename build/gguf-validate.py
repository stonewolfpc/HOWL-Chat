"""Structural check for GGUF files that have no reference copy to compare to.

Size alone cannot tell a complete 7.5 GB quant from a truncated one. This parses
the real GGUF layout — magic, version, tensor count, metadata pairs, then the
tensor-info table — and asks the only question that matters: do the tensors it
declares actually fit in the bytes on disk?

  COMPLETE      every tensor's [offset, offset+size) lands inside the file
  SHORT         declared tensors run past end-of-file (truncated download/paste)
  BAD-HEADER    magic/version/counts are not a GGUF header at all

Byte accounting uses the ggml block-size table. A tensor type outside that table
falls back to offset monotonicity, and the file is reported as PARTIAL rather
than COMPLETE so an unknown quant is never silently blessed.

Usage: python build/gguf-validate.py FILE [FILE ...]
"""
import io
import os
import struct
import sys

# ggml type -> (block elements, bytes per block)
GGML_TYPES = {
    0: ('F32', 1, 4), 1: ('F16', 1, 2), 2: ('Q4_0', 32, 18), 3: ('Q4_1', 32, 20),
    6: ('Q5_0', 32, 22), 7: ('Q5_1', 32, 24), 8: ('Q8_0', 32, 34), 9: ('Q8_1', 32, 40),
    10: ('Q2_K', 256, 84), 11: ('Q3_K', 256, 110), 12: ('Q4_K', 256, 144),
    13: ('Q5_K', 256, 176), 14: ('Q6_K', 256, 210), 15: ('Q8_K', 256, 292),
    16: ('IQ2_XXS', 256, 66), 17: ('IQ2_XS', 256, 74), 18: ('IQ3_XXS', 256, 98),
    19: ('IQ1_S', 256, 50), 20: ('IQ4_NL', 32, 18), 21: ('IQ3_S', 256, 110),
    22: ('IQ2_S', 256, 82), 23: ('IQ4_XS', 256, 136), 24: ('I8', 1, 1),
    25: ('I16', 1, 2), 26: ('I32', 1, 4), 27: ('I64', 1, 8), 28: ('F64', 1, 8),
    29: ('IQ1_M', 256, 56), 30: ('BF16', 1, 2), 34: ('TQ1_0', 256, 54),
    35: ('TQ2_0', 256, 66), 39: ('MXFP4', 32, 17),
}

SCALAR = {0: 1, 1: 1, 2: 2, 3: 2, 4: 4, 5: 4, 6: 4, 7: 1, 10: 8, 11: 8, 12: 8}


class Reader:
    def __init__(self, f, size):
        self.f, self.size, self.pos = f, size, 0

    def take(self, n):
        if self.pos + n > self.size:
            raise ValueError('read past end of file at %d' % self.pos)
        self.f.seek(self.pos)
        b = self.f.read(n)
        self.pos += n
        return b

    def u32(self):
        return struct.unpack('<I', self.take(4))[0]

    def u64(self):
        return struct.unpack('<Q', self.take(8))[0]

    def string(self):
        n = self.u64()
        if n > self.size:
            raise ValueError('string length %d exceeds file' % n)
        return self.take(n)

    def skip_value(self, vtype):
        if vtype in SCALAR:
            self.take(SCALAR[vtype])
        elif vtype == 8:
            self.string()
        elif vtype == 9:
            inner = self.u32()
            count = self.u64()
            if count > self.size:
                raise ValueError('array count %d exceeds file' % count)
            if inner == 8:
                for _ in range(count):
                    self.string()
            elif inner in SCALAR:
                self.take(SCALAR[inner] * count)
            else:
                raise ValueError('array of unknown type %d' % inner)
        else:
            raise ValueError('unknown metadata type %d' % vtype)


def check(path):
    size = os.path.getsize(path)
    with io.open(path, 'rb') as f:
        r = Reader(f, size)
        try:
            magic = r.take(4)
        except ValueError:
            return 'SHORT', 'file is %d bytes, too small to hold a GGUF header' % size
        if magic != b'GGUF':
            return 'BAD-HEADER', 'magic is %r' % magic
        version = r.u32()
        if version not in (1, 2, 3):
            return 'BAD-HEADER', 'version %d' % version
        tensor_count = r.u64()
        kv_count = r.u64()
        if tensor_count > 100000 or kv_count > 100000:
            return 'BAD-HEADER', 'tensor_count=%d kv_count=%d' % (tensor_count, kv_count)

        try:
            alignment = 32
            for _ in range(kv_count):
                key = r.string()
                vtype = r.u32()
                start = r.pos
                r.skip_value(vtype)
                if key == b'general.alignment' and vtype == 4:
                    # Read the value back from where it started: skip_value has
                    # already walked past it.
                    f.seek(start)
                    alignment = struct.unpack('<I', f.read(4))[0] or alignment

            dims_seen = 0
            tensors = []
            for _ in range(tensor_count):
                r.string()
                n_dims = r.u32()
                if n_dims > 8:
                    return 'BAD-HEADER', 'tensor with %d dims' % n_dims
                elements = 1
                for _ in range(n_dims):
                    elements *= r.u64()
                ttype = r.u32()
                offset = r.u64()
                tensors.append((elements, ttype, offset))
                dims_seen += 1
        except ValueError as e:
            return 'SHORT', str(e)

        data_start = r.pos
        if alignment <= 1:
            alignment = 32
        data_start += (-data_start) % alignment
        data_size = size - data_start
        if data_size <= 0:
            return 'SHORT', 'no data section (header ends at %d of %d)' % (data_start, size)

        unknown = set()
        worst_end = 0
        max_offset = 0
        for elements, ttype, offset in tensors:
            max_offset = max(max_offset, offset)
            if ttype not in GGML_TYPES:
                unknown.add(ttype)
                continue
            _, block, per_block = GGML_TYPES[ttype]
            blocks = (elements + block - 1) // block
            worst_end = max(worst_end, offset + blocks * per_block)

        if unknown:
            fit = max_offset < data_size
            return ('PARTIAL' if fit else 'SHORT'), \
                'unknown types %s; deepest offset %d vs data %d' % (
                    sorted(unknown), max_offset, data_size)

        if worst_end > data_size:
            return 'SHORT', 'tensors need %d bytes, file has %d (short by %d)' % (
                worst_end, data_size, worst_end - data_size)
        slack = data_size - worst_end
        if slack >= alignment:
            return 'PARTIAL', 'tensors end %d bytes before EOF' % slack
        return 'COMPLETE', '%d tensors, %d dims, data %d bytes' % (
            dims_seen, dims_seen, data_size)


def main(paths):
    if not paths:
        print('usage: python build/gguf-validate.py FILE [FILE ...]')
        return 2
    bad = 0
    for p in paths:
        try:
            verdict, detail = check(p)
        except OSError as e:
            verdict, detail = 'ERROR', str(e)
        if verdict != 'COMPLETE':
            bad += 1
        print('  %-10s %s  (%s)' % (verdict, os.path.basename(p), detail))
    print('SUMMARY %d files, %d not COMPLETE' % (len(paths), bad))
    return 0


if __name__ == '__main__':
    sys.exit(main(sys.argv[1:]))
