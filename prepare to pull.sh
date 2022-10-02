find . -name "desktop.ini" -exec rm {} \;
echo 'done'
read -rsp $'Press any key to continue...\n' -n 1 key