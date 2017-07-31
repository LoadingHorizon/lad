#!/usr/bin/env bash

WORK_ROOT=`pwd`
export GOPATH=${WORK_ROOT}
export LD_LIBRARY_PATH="${WORK_ROOT}/third/protobuf-3.2.0/lib/"
export PATH="${WORK_ROOT}/bin:${PATH}"
alias protoc="${WORK_ROOT}/third/protobuf-3.2.0/bin/protoc"

function warn() {
    echo "[31m$1[0m"
}

function info() {
    echo "[32m$1[0m"
}

function install_protoc() {
    output=${WORK_ROOT}/third/protobuf-3.2.0
    [[ -e ${output} ]] && info "${output} is exists" && return 0
    mkdir -p third && cd third &&
    wget https://github.com/google/protobuf/releases/download/v3.2.0/protobuf-cpp-3.2.0.tar.gz &&
    tar -xzvf protobuf-cpp-3.2.0.tar.gz && cd protobuf-3.2.0 &&
    ./configure --prefix=${output} && make -j8 && make install &&
    info "install protoc success" && return 0
    warn "install protoc fail" && return 255
}

function install_go_dependency() {
    # go get -v -u github.com/emirpasic/gods
    go get -v -u github.com/siddontang/go-mysql
    go get -v -u github.com/siddontang/go
    go get -v -u golang.org/x/net
    go get -v -u github.com/BurntSushi/toml &&
    go get -v -u github.com/go-sql-driver/mysql &&
    go get -v -u github.com/jmoiron/sqlx &&
    go get -v -u github.com/juju/errors &&
    go get -v -u github.com/ngaut/log &&
    go get -v -u github.com/pingcap/check &&
    go get -v -u github.com/satori/go.uuid &&
    # go get -v -u github.com/golang/protobuf/proto github.com/golang/protobuf/protoc-gen-go &&
    # go get -v -u google.golang.org/grpc &&
    go get -v -u github.com/cihub/seelog &&
    go get -v -u github.com/spf13/viper &&
    # go get -v -u github.com/satori/go.uuid &&
    # go get -v -u github.com/syndtr/goleveldb/leveldb &&
    # go get -v -u github.com/spaolacci/murmur3 &&
    # go get -v -u gopkg.in/mgo.v2 &&
    # go get -v -u github.com/go-sql-driver/mysql &&
    # go get -v -u github.com/garyburd/redigo/redis &&
    # go get -v -u github.com/chasex/redis-go-cluster &&
    # go get -v -u gopkg.in/gin-gonic/gin.v1 &&
    go get -v -u github.com/BurntSushi/toml &&
    go get -v -u gopkg.in/alecthomas/kingpin.v2 &&
    info "install go dependency success" && return 0
    warn "install go dependency fail" && return 255
}

function codegen() {
    cd proto &&
    protoc --go_out=plugins=grpc:. *.proto &&
    mkdir -p ${WORK_ROOT}/src/rta_proto &&
    mv *.go ${WORK_ROOT}/src/rta_proto &&
    info "gencode success" && return 0
    warn "gencode failed" && return 255
}

function build() {
    for f in $(\ls src/main); do
        if [[ ${f} == 'version.go' ]]; then
            continue
        fi
        go build -ldflags "-X 'main.Version=version:  $(git describe --tags) [git describe --tags]
hashcode: $(git rev-parse HEAD) [git rev-parse HEAD]
datetime: $(date '+%Y-%m-%d %H:%M:%S') [date]
hostname: $(hostname) [hostname]'" src/main/${f} src/main/version.go
        if [[ $? != 0 ]]; then
            warn "build src/main/${f} failed." && return 255
        fi
    done

    ver="release"
    if [[ $1 == "prerelease" || $1 == "offline" || $1 == "release" ]]; then
        ver=$1
    fi

    output=output
    mkdir -p ${output}/conf && mkdir -p ${output}/bin
    cp conf/${ver}/* ${output}/conf

    info "build ${ver} success" && return 0
}

function unit() {
    modules=""
    if [[ ! -z "$1" ]]; then
        modules="$1"
    fi
    echo "module	status	time	coverage	code_linenum	test_linenum	code_gocyclo	test_gocyclo"
    for module in ${modules}; do
        time_coverage=$(go test -cover src/${module}/*.go)
        if [[ $? != 0 ]]; then
            go test src/${module}/*.go
            warn "unit test src/${module} failed." && return 255
        fi
        code_linenum=$(find src/${module} -name "*.go" | grep -v _test.go | xargs wc -l | grep -v total | awk '{sum+=$1}END{print sum}')
        test_linenum=$(find src/${module} -name "*.go" | grep _test.go | xargs wc -l | grep -v total | awk '{sum+=$1}END{print sum}')
        code_gocyclo=$(gocyclo $(find src/${module} -name "*.go" | grep -v _test.go) | awk '{sum+=$1}END{print sum}')
        test_gocyclo=$(gocyclo $(find src/${module} -name "*.go" | grep _test.go) | awk '{sum+=$1}END{print sum}')
        echo "${module}\t$(echo ${time_coverage} | awk 'BEGIN{OFS="\t"}{print $1,$3,$5}')\t"\
             "${code_linenum}\t${test_linenum}\t${code_gocyclo}\t${test_gocyclo}"
    done
    info "unit success" && return 0
}

function mstat() {
    gitlab=$(git remote -v | grep origin | grep fetch | awk '{print $2}')
    values=$(echo ${gitlab} | awk -F '[:/.]' 'BEGIN{OFS="\t"}{print $(NF-2),$(NF-1),$0}')
    team=$(echo ${values} | awk '{print $1}')
    proj=$(echo ${values} | awk '{print $2}')
    unit $1 | grep ok | python3 -c "
import sys
import datetime
for line in sys.stdin:
    obj = {}
    values = line.split('\t')
    obj['@team'] = '${team}'
    obj['@proj'] = '${proj}'
    obj['@timestamp'] = datetime.datetime.now() - datetime.timedelta(hours=8)
    obj['module'] = values[0]
    obj['status'] = values[1]
    obj['unittime'] = float(values[2][:-1])
    obj['coverage'] = float(values[3][:-1])
    obj['code_linenum'] = int(values[4])
    obj['test_linenum'] = int(values[5])
    obj['code_gocyclo'] = int(values[6])
    obj['test_gocyclo'] = int(values[7])
    print(('mtech', '{}_{}'.format(obj['@team'], obj['@proj']),
          '{}_{}'.format(obj['module'], obj['@timestamp'].strftime('%Y%m%d')), obj))
    "
}

function usage() {
    info "usage: "
    info "  sh everything.sh <dependency | codegen | unit | build>"
    info "  dependency  安装依赖"
    info "  codegen     生成代码"
    info "  unit        单元测试[]"
    info "  stat        统计"
    info "  build       编译[offline | prerelease | release]"
}

function main() {
    case "$1" in
        "protoc") install_prtoc;;
        "dependency") install_go_dependency;;
        "codegen") codegen;;
        "unit") unit "$2";;
        "stat") mstat "$2";;
        "build") build "$2";;
        *) usage;;
    esac
}

main "$@"
