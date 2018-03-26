rm -rf /tmp/src
cp -r ~/Desktop/src ~/backup
cp -r ~/Desktop/src /tmp
while read p <&3; do
	set $p
	echo $2
	sshpass -pmonkey_pass ssh -t -oStrictHostKeyChecking=no monkey_user@$2 'echo monkey_pass | sudo -S rm -rf ~/Desktop/src; exit'
	echo "copying"
	sshpass -pmonkey_pass scp -oStrictHostKeyChecking=no -r /tmp/src monkey_user@$2:~/Desktop/
	echo "finished"
done 3< /etc/hosts
