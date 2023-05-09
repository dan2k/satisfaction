@echo off
set "TARGET_DIR=D:\dev\satisfaction"
set "TARGET_PATH=%TARGET_DIR%\satisfaction.vbs"
set "SHORTCUT_NAME=satisfaction"
set "STARTUP_FOLDER=%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup"
set "SHORTCUT_PATH=%STARTUP_FOLDER%\%SHORTCUT_NAME%.lnk"
echo Creating shortcut...
echo Set oWS = WScript.CreateObject("WScript.Shell") > CreateShortcut.vbs
echo sLinkFile = "%SHORTCUT_PATH%" >> CreateShortcut.vbs
echo Set oLink = oWS.CreateShortcut(sLinkFile) >> CreateShortcut.vbs
echo oLink.TargetPath = "%TARGET_PATH%" >> CreateShortcut.vbs
echo oLink.WorkingDirectory="%TARGET_DIR%">> CreateShortcut.vbs
echo oLink.Save >> CreateShortcut.vbs
cscript /nologo CreateShortcut.vbs
del CreateShortcut.vbs

echo Shortcut created successfully.
