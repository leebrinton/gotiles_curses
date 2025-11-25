
clone gotileslib:
git clone git@github.com:leebrinton/gotileslib.git 


tell go where to find gotileslib.git:
go mod edit -replace github.com/leebrinton/gotileslib=../gotileslib

go build
