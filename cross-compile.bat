@ECHO OFF

:: Windows
SET GOOS=windows
SET GOARCH=386
go build -o dist/windows/midi2smw_x86.exe

SET GOARCH=amd64
go build -o dist/windows/midi2smw_x64.exe


:: macOS
SET GOOS=darwin
SET GOARCH=amd64
go build -o dist/mac/midi2smw

(
echo --- HOW TO USE ---
echo.
echo Navigate to the directory with the midi2smw executable in it.
echo You'll have to first change the access permissions:
echo 	chmod +x midi2smw
echo.
echo Then execute the program:
echo 	./midi2smw
echo.
echo The first time you try to run it, macOS will pop up a security prompt. Cancel out instead of moving to trash like it suggests.
echo Open up your System Preferences and navigate to Security and Privacy -^> General
echo There will be a message about midi2smw being recently blocked. You may have to click the little lock in the bottom left and input your password to make changes here. Do this, and then click "Open Anyway"
echo.
echo The next time you try to execute the program, you'll probably get one more security warning, but you'll be able to bypass it.
echo From then on, you shouldn't get nagged with warnings any more.
)>"dist/mac/MAC_README.txt"