; Script generated by the Inno Setup Script Wizard.
; SEE THE DOCUMENTATION FOR DETAILS ON CREATING INNO SETUP SCRIPT FILES!

#define MyAppName "EventBucket"
#define MyAppVersion "3.0 Release Preview"
#define MyAppPublisher "EventBucket"
#define MyAppURL "http://www.eventbucket.com.au/"
#define MyAppExeName "EventBucket.exe"

[Setup]
; NOTE: The value of AppId uniquely identifies this application.
; Do not use the same AppId value in installers for other applications.
; (To generate a new GUID, click Tools | Generate GUID inside the IDE.)
AppId={{D515AFDE-322D-49BE-8240-94C5A655BB5C}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
;AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={pf}\{#MyAppName}
DefaultGroupName={#MyAppName}
AllowNoIcons=yes
OutputDir=Z:\inno
LicenseFile=Z:\inno\cc-by-sa-4.0_legalcode.txt
OutputBaseFilename={#MyAppName} {#MyAppVersion} 32bit
SetupIconFile=Z:\EventBucket\EventBucket4.ico
Compression=lzma
SolidCompression=yes

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked
Name: "quicklaunchicon"; Description: "{cm:CreateQuickLaunchIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked; OnlyBelowVersion: 0,6.1

[Files]
Source: "Z:\EventBucket\EventBucket.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "Z:\EventBucket\favicon.ico"; DestDir: "{app}"; Flags: ignoreversion
Source: "Z:\EventBucket\c\*"; DestDir: "{app}"; Flags: ignoreversion
Source: "Z:\EventBucket\h\*"; DestDir: "{app}"; Flags: ignoreversion
Source: "Z:\EventBucket\j\*"; DestDir: "{app}"; Flags: ignoreversion
Source: "Z:\EventBucket\p\*"; DestDir: "{app}"; Flags: ignoreversion
Source: "Z:\EventBucket\v\*"; DestDir: "{app}"; Flags: ignoreversion
; NOTE: Don't use "Flags: ignoreversion" on any shared system files

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{commondesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon
Name: "{userappdata}\Microsoft\Internet Explorer\Quick Launch\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: quicklaunchicon

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#StringChange(MyAppName, '&', '&&')}}"; Flags: nowait postinstall skipifsilent

