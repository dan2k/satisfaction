
'1. create short cut from  satisfaction.vbs 
'2. copy short cut to  startup   example: C:\Users\dan2k\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup
Set WshShell = CreateObject("WScript.Shell" ) 
WshShell.Run "satisfaction.exe", 0 'Must quote command if it has spaces; must escape quotes
Set WshShell = Nothing