while read p <&3; do
	set $p
	echo $2
	sshpass -pseng468 ssh -t -oStrictHostKeyChecking=no  seng468@$2 'echo seng468 | sudo -S adduser monkey_user; echo "monkey_pass" | sudo passwd --stdin monkey_user; sudo usermod -aG wheel monkey_user; sudo groupadd docker; sudo usermod -aG docker monkey_user'
done 3< /etc/hosts
