param(
    [ValidateSet('list','getdefault','setdefault','getvolume','setvolume')]
    [string]$Action = 'getdefault',
    [string]$Match = '',
    [string]$Roles = '0,1,2',
    [int]$Percent = -1
)
$ErrorActionPreference = 'Stop'

$cs = @'
using System;
using System.Runtime.InteropServices;

public static class AudioCtlImpl
{
    [ComImport, Guid("BCDE0395-E52F-467C-8E3D-C4579291692E")]
    class MMDeviceEnumeratorComObject { }

    [ComImport, Guid("870af99c-171d-4f9e-af0d-e63df40c2bc9")]
    class PolicyConfigClient { }

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

    [ComImport, Guid("f8679f50-850a-41cf-9c72-430f290290c8"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
    interface IPolicyConfig
    {
        int GetMixFormat([MarshalAs(UnmanagedType.LPWStr)] string pszDeviceName, out IntPtr ppFormat);
        int GetDeviceFormat([MarshalAs(UnmanagedType.LPWStr)] string pszDeviceName, int bDefault, out IntPtr ppFormat);
        int ResetDeviceFormat([MarshalAs(UnmanagedType.LPWStr)] string pszDeviceName);
        int SetDeviceFormat([MarshalAs(UnmanagedType.LPWStr)] string pszDeviceName, IntPtr pEndpointFormat, IntPtr mixFormat);
        int GetProcessingPeriod([MarshalAs(UnmanagedType.LPWStr)] string pszDeviceName, int bDefault, out IntPtr pmftDefaultPeriod, out IntPtr pmftMinimumPeriod);
        int SetProcessingPeriod([MarshalAs(UnmanagedType.LPWStr)] string pszDeviceName, IntPtr pmftPeriod);
        int GetShareMode([MarshalAs(UnmanagedType.LPWStr)] string pszDeviceName, out IntPtr pMode);
        int SetShareMode([MarshalAs(UnmanagedType.LPWStr)] string pszDeviceName, IntPtr mode);
        int GetPropertyValue([MarshalAs(UnmanagedType.LPWStr)] string pszDeviceName, IntPtr key, out IntPtr pv);
        int SetPropertyValue([MarshalAs(UnmanagedType.LPWStr)] string pszDeviceName, IntPtr key, IntPtr pv);
        int SetDefaultEndpoint([MarshalAs(UnmanagedType.LPWStr)] string pszDeviceName, int role);
        int SetEndpointVisibility([MarshalAs(UnmanagedType.LPWStr)] string pszDeviceName, int bVisible);
    }

    [ComImport, Guid("5CDF2C82-841E-4546-9722-0CF74078229A"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
    interface IAudioEndpointVolume
    {
        int RegisterControlChangeNotify(IntPtr p);
        int UnregisterControlChangeNotify(IntPtr p);
        int GetChannelCount(out uint c);
        int SetMasterVolumeLevel(float level, Guid ctx);
        int SetMasterVolumeLevelScalar(float level, Guid ctx);
        int GetMasterVolumeLevel(out float level);
        int GetMasterVolumeLevelScalar(out float level);
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

    static IMMDeviceEnumerator Enumerator()
    {
        return (IMMDeviceEnumerator)(new MMDeviceEnumeratorComObject());
    }

    public static string[] ListRender()
    {
        var en = Enumerator();
        IMMDeviceCollection coll;
        if (en.EnumAudioEndpoints(EDataFlow.eRender, 0x0000000F, out coll) != 0) return new string[0];
        uint n; coll.GetCount(out n);
        var list = new System.Collections.Generic.List<string>();
        for (uint i = 0; i < n; i++)
        {
            IMMDevice d; coll.Item(i, out d);
            int state; d.GetState(out state);
            string id; d.GetId(out id);
            IPropertyStore st;
            string name = "";
            if (d.OpenPropertyStore(0, out st) == 0)
            {
                PROPVARIANT pv;
                if (st.GetValue(ref PKEY_FriendlyName, out pv) == 0 && pv.vt == 31 && pv.p != IntPtr.Zero)
                    name = Marshal.PtrToStringUni(pv.p);
                Marshal.ReleaseComObject(st);
            }
            list.Add(name + "|state=" + state + "|" + id);
            Marshal.ReleaseComObject(d);
        }
        return list.ToArray();
    }

    public static string GetDefaultRender()
    {
        var en = Enumerator();
        IMMDevice d;
        if (en.GetDefaultAudioEndpoint(EDataFlow.eRender, ERole.eConsole, out d) != 0) return "NONE";
        string id; d.GetId(out id);
        IPropertyStore st; string name = "";
        if (d.OpenPropertyStore(0, out st) == 0)
        {
            PROPVARIANT pv;
            if (st.GetValue(ref PKEY_FriendlyName, out pv) == 0 && pv.vt == 31 && pv.p != IntPtr.Zero)
                name = Marshal.PtrToStringUni(pv.p);
            Marshal.ReleaseComObject(st);
        }
        Marshal.ReleaseComObject(d);
        return name + "|" + id;
    }

    // Sets the default console/multimedia/communications render endpoint whose
    // friendly name contains the given substring. Returns the prior default.
    public static string SetDefaultRender(string match, string rolesCsv)
    {
        string before = GetDefaultRender();
        string targetId = null; string targetName = null;
        foreach (var row in ListRender())
        {
            var parts = row.Split('|');
            if (parts.Length < 3) continue;
            if (!parts[1].EndsWith("=1")) continue; // active only
            if (parts[0].IndexOf(match, StringComparison.OrdinalIgnoreCase) < 0) continue;
            targetName = parts[0]; targetId = parts[2]; break;
        }
        if (targetId == null) throw new Exception("no active render endpoint matching '" + match + "'");
        var pc = (IPolicyConfig)(new PolicyConfigClient());
        foreach (var token in rolesCsv.Split(','))
        {
            int role = int.Parse(token.Trim());
            int hr = pc.SetDefaultEndpoint(targetId, role);
            if (hr != 0) throw new Exception("SetDefaultEndpoint role=" + role + " hr=0x" + hr.ToString("X"));
        }
        System.Threading.Thread.Sleep(400);
        return before + " => " + targetName;
    }

    public static float GetVolume()
    {
        var en = Enumerator();
        IMMDevice d;
        if (en.GetDefaultAudioEndpoint(EDataFlow.eRender, ERole.eConsole, out d) != 0) return -1;
        Guid iid = new Guid("5CDF2C82-841E-4546-9722-0CF74078229A");
        IntPtr p;
        if (d.Activate(ref iid, 0, IntPtr.Zero, out p) != 0) return -1;
        var vol = (IAudioEndpointVolume)Marshal.GetObjectForIUnknown(p);
        float v; vol.GetMasterVolumeLevelScalar(out v);
        Marshal.ReleaseComObject(vol); Marshal.Release(p); Marshal.ReleaseComObject(d);
        return v;
    }

    public static void SetVolume(float scalar)
    {
        var en = Enumerator();
        IMMDevice d;
        if (en.GetDefaultAudioEndpoint(EDataFlow.eRender, ERole.eConsole, out d) != 0) throw new Exception("no default render");
        Guid iid = new Guid("5CDF2C82-841E-4546-9722-0CF74078229A");
        IntPtr p;
        if (d.Activate(ref iid, 0, IntPtr.Zero, out p) != 0) throw new Exception("activate volume failed");
        var vol = (IAudioEndpointVolume)Marshal.GetObjectForIUnknown(p);
        vol.SetMasterVolumeLevelScalar(scalar, Guid.Empty);
        Marshal.ReleaseComObject(vol); Marshal.Release(p); Marshal.ReleaseComObject(d);
    }
}
'@

Add-Type -TypeDefinition $cs -Language CSharp

switch ($Action) {
    'list' {
        [AudioCtlImpl]::ListRender() | ForEach-Object { Write-Output "  $_" }
    }
    'getdefault' {
        Write-Output ("DEFAULT_RENDER=" + [AudioCtlImpl]::GetDefaultRender())
    }
    'setdefault' {
        if ($Match -eq '') { throw '-Match is required for setdefault' }
        Write-Output ("SWITCH " + [AudioCtlImpl]::SetDefaultRender($Match, $Roles))
    }
    'getvolume' {
        $v = [AudioCtlImpl]::GetVolume()
        Write-Output ("VOLUME={0}%" -f [int]([math]::Round($v * 100)))
    }
    'setvolume' {
        if ($Percent -lt 0) { throw '-Percent is required for setvolume' }
        [AudioCtlImpl]::SetVolume([float]($Percent / 100.0))
        $v = [AudioCtlImpl]::GetVolume()
        Write-Output ("VOLUME={0}%" -f [int]([math]::Round($v * 100)))
    }
}
