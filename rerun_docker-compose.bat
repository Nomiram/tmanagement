docker-compose down
docker-compose build
@IF errorlevel 1 echo ERROR:%ERRORLEVEL% && pause && exit
REM start /B docker-compose up 
docker-compose up 
@IF errorlevel 1 echo ERROR:%ERRORLEVEL% && pause

