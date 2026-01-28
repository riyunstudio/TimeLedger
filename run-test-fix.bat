@echo off
cd /d d:\project\TimeLedger
go test ./testing/test -run "TestSmartMatchingService_InviteTalent|TestSmartMatchingService_GetTalentStats" -v
pause
