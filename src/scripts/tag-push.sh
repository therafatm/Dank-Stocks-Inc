
images=$(docker images "src*" |  awk '{ if (NR > 1) print $1}')
regip=5111

for image in $images
do
	docker tag $image localhost:$regip/$image
	docker push localhost:$regip/$image
done