Set oShell = CreateObject ("Wscript.Shell")
Dim strArgs
strArgs = "cmd.exe /c venv\Scripts\activate.bat && pythonw.exe pcdeck.py"
oShell.Run strArgs, 0, false
