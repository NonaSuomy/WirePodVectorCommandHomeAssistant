# Wire-Pod Vector Command HomeAssistant
Wire-Pod GO plugin for Anki Vector Robot to talk to and send commands to Home Assistant via the Conversation API

![vector-vector-robot](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/8ca0d3bf-df97-4cc2-b1a1-89dbe5d1865e)

https://youtube.com/shorts/i7WPcnAWji8

## Docker ##
I use docker to run Vector with Wire-Pod my configuration looks like this:

```docker
FROM ubuntu:latest

# Edit your timezone here:
ENV TZ=Europe/London
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt-get update \
 && apt-get install -y sudo

RUN adduser --disabled-password --gecos '' wirepod
RUN adduser wirepod sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

USER wirepod

RUN sudo apt-get update && sudo apt-get install -y git
RUN sudo mkdir /wire-pod
RUN sudo chown -R wirepod:wirepod /wire-pod

RUN cd /wire-pod
RUN git clone https://github.com/kercre123/wire-pod/ wire-pod

WORKDIR /wire-pod

RUN sudo STT=vosk /wire-pod/setup.sh

WORKDIR /wire-pod/chipper

CMD sudo /wire-pod/chipper/start.sh
```

```bash
sudo docker build -t wire-pod .
```

My Docker Compose looks like this
```yaml
# Wire-Pod Container for Docker
# https://github.com/kercre123/wire-pod/wiki/Installation
version: "3.9"
services:
  wirepod:
    container_name: "wirepod"
    image: "wire-pod:latest"
    #labels:
    #  - "com.centurylinklabs.watchtower.enable=true"
    #  - "com.centurylinklabs.watchtower.monitor-only=false"
    volumes:
      - /var/lib/libvirt/images/usb001/docker-storage/wirepod/data:/data
    # - /run/dbus:/run/dbus:ro # This should technically get dbus running for Avahi but didn't seem to work for me.
                               # May have to use dbus-broker instead of dbus-daemon.
    mac_address: "9e:de:ad:be:ef:42"
    hostname: escapepod
    #networks:
    #  - dhcp
    ports:
      - 80:80
      - 8080:8080
      - 443:443
      - 8084:8084
      - 5353:5353/udp
    restart: "no" # Turn this to unless-stopped if you think everything is ok 
                        # If you use docker-net-dhcp like me keep this set to no so your dockerd doesn't hang on boot. 
#networks:
#  dhcp:
#    name: "dbrv42"
#    external: true
```
```bash
docker compose up -d
```

## Compiling ##

After that attach to the docker

```bash
docker exec -it wirepod /bin/bash
```
You should be met with the cli prompt
```bash
wirepod@escapepod:/wire-pod/chipper$
cd plugins
mkdir commandha
cd commandha
wget https://raw.githubusercontent.com/NonaSuomy/WirePodVectorCommandHomeAssistant/main/commandha.go
```
Edit the commandha.go file to add your long-lived access token and the IP of your HA instance

```
url := "http://HAIPHERE:8123/api/conversation/process" // Replace with your Home Assistant IP
token := "LONGTOKENHERE" // Replace with your token
```

Compile the GO plugin to the root directory of /wire-pod/chipper/plugins 
```bash
sudo /usr/local/go/bin/go build -buildmode=plugin -o /wire-pod/chipper/plugins/commandha.so commandha.go
```

Restart your wirepod docker

```bash
docker container restart wirepod
```
## Testing, testing, 123 ##

In the console log you should see

```
docker logs wirepod
```

```
SDK info path: /tmp/.anki_vector/
API config successfully read
Loaded jdocs file
Loaded bot info file, known bots: [00######]
Reading session certs for robot IDs
Loaded 54 intents and 54 matches (language: en-US)
Initiating vosk voice processor with language en-US
Opening VOSK model (../vosk/models/en-US/model)
Initializing VOSK recognizers
VOSK initiated successfully
```
After this it should load our compiled plugin
```
Loading plugins
Loading plugin: comandha.so
Utterances []string in plugin comandha.so are OK
Action func in plugin comandha.so is OK
Name string in plugin Home Assistant Control is OK
comandha.so loaded successfully
```
Then continue on...
```
Starting webserver at port 8080 (http://localhost:8080)
Starting jdocs pinger ticker
Starting SDK app
Starting server at port 80 for connCheck
Initiating TLS listener, cmux, gRPC handler, and REST handler
Configuration page: http://VECTORIP:8080
Registering escapepod.local on network (loop)
Starting chipper server at port 443
Starting chipper server at port 8084 for 2.0.1 compatibility
wire-pod started successfully!
Jdocs: Incoming ReadDocs request, Robot ID: vic:00######, Item(s) to return: 
[doc_name:"vic.AppTokens"]
Successfully got jdocs from 00######
Vector discovered on network, broadcasting mDNS
Broadcasting mDNS now (outside of timer loop)
```
Then a successful detection and command looks like this ("Hey Vector" pause... "assist turn off the family room lights")
```
This is a custom intent or plugin!
Bot 00###### request served.
Bot 00###### Stream type: OPUS
(Bot 00######, Vosk) Processing...
Using general recognizer
(Bot 00######) End of speech detected.
Bot 00###### Transcribed text: assist turn off the family room lights
Bot 00###### matched plugin Home Assistant Control, executing function
Bot 00###### plugin Home Assistant Control, response Turned off the lights
```
Then turn the back on...
```
This is a custom intent or plugin!
Bot 00###### request served.
Bot 00###### Stream type: OPUS
(Bot 00######, Vosk) Processing...
Using general recognizer
(Bot 00######) End of speech detected.
Bot 00###### Transcribed text: assist turn on the family room lights
Bot 00###### matched plugin Home Assistant Control, executing function
Bot 00###### plugin Home Assistant Control, response Turned on the lights
This is a custom intent or plugin!
Bot 00###### request served.
```

![200w](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/212e3236-4ec2-44a8-bb3a-dc37d666fb45)

Carry on...
```
Haven't recieved a conn check from 00###### in 15 seconds, will ping jdocs on next check
Broadcasting mDNS now (outside of timer loop)
Jdocs: Incoming ReadDocs request, Robot ID: vic:00######, Item(s) to return: 
[doc_name:"vic.AppTokens"]
Successfully got jdocs from 00######
```
A failure speech detection should look something like this
```
(Bot 00######, Vosk) Processing...
Using general recognizer
(Bot 00######) End of speech detected.
Bot 00###### Transcribed text: a third
Not a custom intent
No intent was matched.
Bot 00###### Intent Sent: intent_system_noaudio
No Parameters Sent
Bot 00###### Stream type: OPUS
(Bot 00######, Vosk) Processing...
Using general recognizer
(Bot 00######) End of speech detected.
Bot 00###### Transcribed text: isis to turn off the living room like
Not a custom intent
No intent was matched.
Bot 00###### Intent Sent: intent_system_noaudio
No Parameters Sent
Bot 00###### Stream type: OPUS
(Bot 00######, Vosk) Processing...
Using general recognizer
(Bot 00######) End of speech detected.
Bot 00###### Transcribed text: assist turn off the family room like
Bot 00###### matched plugin Home Assistant Control, executing function
Bot 00###### plugin Home Assistant Control, response Sorry, I couldn't understand that
This is a custom intent or plugin!
Bot 00###### request served.
Bot 00###### Stream type: OPUS
(Bot 00######, Vosk) Processing...
Using general recognizer
(Bot 00######) End of speech detected.
Bot 00###### Transcribed text: assist
Bot 00###### matched plugin Home Assistant Control, executing function
Bot 00###### plugin Home Assistant Control, response Sorry, I couldn't understand that
```

## Useful links ##

Wire-Pod Go compile information:
https://github.com/kercre123/wire-pod/wiki/For-Developers-and-Tinkerers

My issue trying to get wire-pod working in docker:
https://github.com/kercre123/wire-pod/issues/201

Looks like RGarrett93 is priming the ability to have wire-pod as an add-on to HA woo hoo! Help them if you can! I hope what I did does as well!
[https://community.home-assistant.io/t/anki-vector-integration/](https://community.home-assistant.io/t/anki-vector-integration/85156/33)

![giphy](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/f97d13ac-f585-4c27-aa5f-0e2ed5e7d675)

**Note:** _I have no idea how I was able to accomplish this! <insert codepilot conversation here> Please feel free to do a PR and make this better if you want and let me know what you did! I just attempted to break ground for others as I'm a basic hacker at best and just really wanted to issue commands to HA through Vector locally instead of using Alexa on him..._
