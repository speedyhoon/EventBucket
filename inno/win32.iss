#define MyAppName "EventBucket"
#define MyAppVersion "3.04"
#define MyAppURL "http://www.eventbucket.com.au/"
#define MyAppExeName "EventBucket.exe"
#define Z "\\camtop\EventBucket"

[Setup]
; NOTE: The value of AppId uniquely identifies this application.
; Do not use the same AppId value in installers for other applications.
; (To generate a new GUID, click Tools | Generate GUID inside the IDE.)
AppId={{D515AFDE-322D-49BE-8240-94C5A655BB5C}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppName}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={pf}\{#MyAppName}
DefaultGroupName={#MyAppName}
AllowNoIcons=yes
OutputDir={#Z}\inno
LicenseFile={#Z}\license
OutputBaseFilename={#MyAppName} {#MyAppVersion} 32bit
SetupIconFile={#Z}\icon\app.ico
Compression=lzma
SolidCompression=yes

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"
Name: "quicklaunchicon"; Description: "{cm:CreateQuickLaunchIcon}"; GroupDescription: "{cm:AdditionalIcons}"; OnlyBelowVersion: 0,6.1

[Files]
Source: "{#Z}\EventBucket\EventBucket.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#Z}\icon\favicon.ico"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#Z}\EventBucket\c\*"; DestDir: "{app}\c"; Flags: ignoreversion
Source: "{#Z}\EventBucket\h\*"; DestDir: "{app}\h"; Flags: ignoreversion
Source: "{#Z}\EventBucket\j\*"; DestDir: "{app}\j"; Flags: ignoreversion
Source: "{#Z}\EventBucket\v\*"; DestDir: "{app}\v"; Flags: ignoreversion
Source: "{#Z}\EventBucket\w\*"; DestDir: "{app}\w"; Flags: ignoreversion
Source: "{#Z}\license"; DestDir: "{app}"; Flags: ignoreversion
; NOTE: Don't use "Flags: ignoreversion" on any shared system files

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{group}\{#MyAppName} dark"; Filename: "{app}\{#MyAppExeName}"; Parameters: "-dark"
Name: "{commondesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon
Name: "{commondesktop}\{#MyAppName} dark"; Filename: "{app}\{#MyAppExeName}"; Parameters: "-dark"; Tasks: desktopicon
Name: "{userappdata}\Microsoft\Internet Explorer\Quick Launch\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: quicklaunchicon

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#StringChange(MyAppName, '&', '&&')}}"; Flags: nowait postinstall skipifsilent