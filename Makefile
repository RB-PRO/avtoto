all: run

run:
	go run main.go GetBrandsByCode.go SearchGetParts2.go SearchStart.go httpCustom.go

push:
	git push git@github.com:RB-PRO/avtoto.git

pull:
	git pull git@github.com:RB-PRO/avtoto.git

pushW:
	git push https://github.com/RB-PRO/avtoto.git

pullW:
	git pull https://github.com/RB-PRO/avtoto.git
