$ErrorActionPreference = 'Stop'

$cs = @'
using System;
using System.Runtime.InteropServices;

public static class AudioEP
{
    [ComImport, Guid("BCDE0395-E52F-467C-8E3D-C4579291692E")]
    class MMDeviceEnumeratorComObject { }

    public enum EDataFlow { eRender = 0, eCapture = 1, eAll = 2 }
    public enum ERole { eConsole = 0, eMultimedia = 1, eCommunications = 2 }

    [ComImport, Guid("A95664D2-9614-4F35-A746-DE8DB63617E6"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
    interface IMMDeviceEnumerator
    {
        int EnumAudioEndpoints(EDataFlow dataFlow, int stateMask, out IMMDeviceCollection ppDevices);
        int GetDefaultAudioEndpoint(EDataFlow dataFlow, ERole role, out IMMDevice ppDevice);
        int GetDevice([MarshalAs(UnmanagedType.LPWStr)] string pwstrId, out IMMDevice ppDevice);
        int RegisterEndpointNotificationCallback(IntPtr pClient);
        int UnregisterEndpointNotificationCallback(IntPtr pClient);
    }

    [ComImport, Guid("0BD7A1BE-7A1A-44DB-8397-CC5392387B5E"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
    interface IMMDeviceCollection
    {
        int GetCount(out uint pcDevices);
        int Item(uint nDevice, out IMMDevice ppDevice);
    }

    [ComImport, Guid("D666063F-1587-4E43-81F1-B948E807363F"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
    interface IMMDevice
    {
        int Activate(ref Guid iid, int dwClsCtx, IntPtr pActivationParams, out IntPtr ppInterface);
        int OpenPropertyStore(int stgmAccess, out IPropertyStore ppProperties);
        int GetId([MarshalAs(UnmanagedType.LPWStr)] out string ppstrId);
        int GetState(out int pdwState);
    }

    [ComImport, Guid("886d8eeb-8cf2-4446-8d02-cdba1dbdcf99"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
    interface IPropertyStore
    {
        int GetCount(out uint cProps);
        int GetAt(uint iProp, out PROPERTYKEY pkey);
        int GetValue(ref PROPERTYKEY key, out PROPVARIANT pv);
        int SetValue(ref PROPERTYKEY key, ref PROPVARIANT pv);
        int Commit();
    }

    [StructLayout(LayoutKind.Sequential)]
    public struct PROPERTYKEY { public Guid fmtid; public uint pid; }

    [StructLayout(LayoutKind.Sequential)]
    public struct PROPVARIANT
    {
        public ushort vt; public ushort r1; public ushort r2; public ushort r3;
        public IntPtr p; public IntPtr p2;
    }

    static PROPERTYKEY PKEY_FriendlyName = new PROPERTYKEY {
        fmtid = new Guid("a45c254e-df1c-4efd-8020-67d146a850e0"), pid = 14 };

    static PROPERTYKEY PKEY_DeviceDesc = new PROPERTYKEY {
        fmtid = new Guid("a45c254e-df1c-4efd-8020-67d146a850e0"), pid = 2 };

    static string ReadProp(IPropertyStore store, PROPERTYKEY key)
    {
        PROPVARIANT pv;
        int hr = store.GetValue(ref key, out pv);
        if (hr != 0) return "";
        if (pv.vt == 31 && pv.p != IntPtr.Zero) return Marshal.PtrToStringUni(pv.p);
        return "";
    }

    static string Describe(IMMDevice dev)
    {
        string id; dev.GetId(out id);
        IPropertyStore store;
        string name = "", desc = "";
        if (dev.OpenPropertyStore(0, out store) == 0)
        {
            name = ReadProp(store, PKEY_FriendlyName);
            desc = ReadProp(store, PKEY_DeviceDesc);
            Marshal.ReleaseComObject(store);
        }
        int state; dev.GetState(out state);
        return name + "|" + desc + "|state=" + state + "|" + id;
    }

    public static string[] List(int flow)
    {
        var en = (IMMDeviceEnumerator)(new MMDeviceEnumeratorComObject());
        IMMDeviceCollection coll;
        int hr = en.EnumAudioEndpoints((EDataFlow)flow, 0x0000000F, out coll);
        if (hr != 0) return new string[] { "ERR EnumAudioEndpoints hr=0x" + hr.ToString("X") };
        uint n; coll.GetCount(out n);
        var list = new System.Collections.Generic.List<string>();
        for (uint i = 0; i < n; i++)
        {
            IMMDevice d; coll.Item(i, out d);
            list.Add(Describe(d));
            Marshal.ReleaseComObject(d);
        }
        Marshal.ReleaseComObject(coll);
        Marshal.ReleaseComObject(en);
        return list.ToArray();
    }

    public static string Default(int flow, int role)
    {
        var en = (IMMDeviceEnumerator)(new MMDeviceEnumeratorComObject());
        IMMDevice d;
        int hr = en.GetDefaultAudioEndpoint((EDataFlow)flow, (ERole)role, out d);
        if (hr != 0) return "NONE hr=0x" + hr.ToString("X");
        string s = Describe(d);
        Marshal.ReleaseComObject(d);
        Marshal.ReleaseComObject(en);
        return s;
    }
}
'@

Add-Type -TypeDefinition $cs -Language CSharp

Write-Output "=== CAPTURE ENDPOINTS (state mask: active|disabled|notpresent|unplugged) ==="
[AudioEP]::List(1) | ForEach-Object { Write-Output "  $_" }
Write-Output ""
Write-Output "=== RENDER ENDPOINTS ==="
[AudioEP]::List(0) | ForEach-Object { Write-Output "  $_" }
Write-Output ""
Write-Output "=== DEFAULTS ==="
Write-Output ("  capture/console : " + [AudioEP]::Default(1,0))
Write-Output ("  capture/comms   : " + [AudioEP]::Default(1,2))
Write-Output ("  render/console  : " + [AudioEP]::Default(0,0))
Write-Output ("  render/comms    : " + [AudioEP]::Default(0,2))
