docker-compose down
docker-compose build
@IF errorlevel 1 echo ERROR:%ERRORLEVEL% && pause
start /B docker-compose up
@IF errorlevel 1 echo ERROR:%ERRORLEVEL% && pause

