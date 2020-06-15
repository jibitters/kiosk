#!/usr/bin/env sh

go vet ./...
golint -set_exit_status=1 ./...

case $DOCKER in
	"")
    ginkgo -r -p --nodes=8 --v --trace -race
	  ;;

	*)
		ginkgo -r -p --nodes=8 --v --trace -race -- --pg.host "$DOCKER"
		;;
esac
