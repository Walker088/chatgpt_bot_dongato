
dockerRun() {
    docker run -d -it -v "$(pwd)/.project.env:/.project.env" --name dongato cloud.canister.io:5000/walker088/chatgpt_bot_dongato
}
dockerStop() {
    docker stop dongato
}
dockerRemove() {
    docker rm dongato
}
dockerRestart() {
    dockerStop
    dockerRemove
    dockerRun
}
dockerBuild() {
    docker build . -t cloud.canister.io:5000/walker088/chatgpt_bot_dongato
}

if [[ $# -eq 0 ]] ; then
    echo 'Please provide one of the arguments (e.g., bash flyway-migrate.sh info):
    1 > run
    2 > stop
    3 > remove
    4 > restart
    5 > build'

elif [[ $1 == run ]]; then
    dockerRun

elif [[ $1 == stop ]]; then
    dockerStop

elif [[ $1 == remove ]]; then
    dockerRemove

elif [[ $1 == restart ]]; then
    dockerRestart

elif [[ $1 == build ]]; then
    dockerBuild

fi
