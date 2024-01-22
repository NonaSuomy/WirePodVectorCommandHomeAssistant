# Wire-Pod Vector Command HomeAssistant
Wire-Pod GO plugin for Anki Vector Robot to talk to and send commands to Home Assistant via the Conversation API

![vector-vector-robot](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/8ca0d3bf-df97-4cc2-b1a1-89dbe5d1865e)

## YouTube Short ##

[![YouTube Short](https://img.youtube.com/vi/i7WPcnAWji8/0.jpg)](https://www.youtube.com/watch?v=i7WPcnAWji8)

^ https://youtube.com/shorts/i7WPcnAWji8 ^

## Home Assistant Container for Docker ##

```yaml
#HomeAssistant Container for Docker
version: "3.9"
services:
  homeassistant:
    container_name: homeassistant
    image: "ghcr.io/home-assistant/home-assistant:stable"
    volumes:
      - /var/lib/libvirt/images/usb001/docker-storage/homeassistant/config:/config
      - /etc/localtime:/etc/localtime:ro
    mac_address: de:ad:be:ef:00:08
    networks:
      - dhcp
    devices:
      - /dev/ttyACM0:/dev/ttyACM0
    restart: "no"
    privileged: true

networks:
  dhcp:
    name: dbrv100
    external: true
```

## Home-LLM on a GPU server ##
[![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/a4cd755e-3d3f-435d-875a-e3eeb0b3e420)](https://my.home-assistant.io/redirect/hacs_repository/?category=Integration&repository=home-llm&owner=acon96)

Clone or download the repository https://github.com/oobabooga/text-generation-webui install it however you please I just run it on the metal of a GPU server.

**Note:** _You currently need an unreleased tag of text-generation-webui to work with the V2 model of Home-LLM which you can get like this until fully released_

```bash
git clone https://github.com/oobabooga/text-generation-webui.git
cd text-generation-webui
git checkout -b llamacpp_0.2.29 origin/llamacpp_0.2.29
git branch
git pull
echo "--listen --api --model Home-3B-v2.q8_0.gguf --n-gpu-layers 33" >> CMD_FLAGS.txt
./start_linux.sh
```

Edit CMD_FLAGS.txt before installing uncomment the line `# --listen --api` (Remove the number sign and space) I also added `--model Home-3B-v2.q8_0.gguf --n-gpu-layers 33` for my GPU to start the model on boot.

Run the start_linux.sh, start_windows.bat, start_macos.sh, or start_wsl.bat script depending on your OS.

Select your GPU vendor when asked.

Once the installation ends, browse to http://GPUServerIP:7860/?__theme=dark.

Download the gguf file from acon96/Home-3B-v2-GGUF and Home-3B-v2.q8_0.gguf for max GPU potential in the text-generation-webui.

Select the reload button when it is done downloading (2.8G) the file Home-3B-v2.q8_0.gguf should show up in the Model list then hit Load and Save settings.

![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/8899f18c-a6af-4ba8-87b2-ed9719de7820)

Excerpt from Home-LLM
```
Performance of running the model on a Raspberry Pi

The RPI4 4GB that I have was sitting right at 1.5 tokens/sec for prompt eval and 1.6 tokens/sec for token generation when running the Q4_K_M quant. I was reliably getting responses in 30-60 seconds after the initial prompt processing which took almost 5 minutes. It depends significantly on the number of devices that have been exposed as well as how many states have changed since the last invocation because of llama.cpp caches KV values for identical prompt prefixes.

It is highly recommended to set up text-generation-webui on a separate machine that can take advantage of a GPU.
```

Start-up on boot of the GPU server.
```
#/etc/systemd/system/textgen.service
[Unit]
After=network.target

[Service]
Type=simple
ExecStart=/bin/bash /opt/code/text-generation-webui/start_linux.sh
User=textgen
Group=textgen

[Install]
WantedBy=multi-user.target
```

Add user textgen and to groups video/render

```bash
sudo useradd -m -s /bin/bash -G render,video textgen
sudo passwd textgen
id textgen
uid=1001(textgen) gid=1001(textgen) groups=1001(textgen),44(video),109(render)
sudo chown -R textgen:textgen /opt/text-generation-webui/
```

Test your user to make sure it can run text-generation-webui

```bash
su textgen
cd /opt/text-generation-webui
./start_linux.sh
08:11:51-750558 INFO     Starting Text generation web UI
08:11:51-753175 WARNING
                         You are potentially exposing the web UI to the entire internet without any access password.
                         You can create one with the "--gradio-auth" flag like this:

                         --gradio-auth username:password

                         Make sure to replace username:password with your own.
08:11:51-822630 INFO     Loading Home-3B-v2.q8_0.gguf
08:11:51-907601 INFO     llama.cpp weights detected: models/Home-3B-v2.q8_0.gguf
ggml_init_cublas: GGML_CUDA_FORCE_MMQ:   no
ggml_init_cublas: CUDA_USE_TENSOR_CORES: yes
ggml_init_cublas: found 1 ROCm devices:
  Device 0: Radeon RX Vega, compute capability 9.0, VMM: no
llama_model_loader: loaded meta data with 19 key-value pairs and 453 tensors from models/Home-3B-v2.q8_0.gguf (version GGUF V3 (latest))
llama_model_loader: Dumping metadata keys/values. Note: KV overrides do not apply in this output.
llama_model_loader: - kv   0:                       general.architecture str              = phi2
llama_model_loader: - kv   1:                               general.name str              = Phi2
llama_model_loader: - kv   2:                        phi2.context_length u32              = 2048
llama_model_loader: - kv   3:                      phi2.embedding_length u32              = 2560
llama_model_loader: - kv   4:                   phi2.feed_forward_length u32              = 10240
llama_model_loader: - kv   5:                           phi2.block_count u32              = 32
llama_model_loader: - kv   6:                  phi2.attention.head_count u32              = 32
llama_model_loader: - kv   7:               phi2.attention.head_count_kv u32              = 32
llama_model_loader: - kv   8:          phi2.attention.layer_norm_epsilon f32              = 0.000010
llama_model_loader: - kv   9:                  phi2.rope.dimension_count u32              = 32
llama_model_loader: - kv  10:                          general.file_type u32              = 7
llama_model_loader: - kv  11:               tokenizer.ggml.add_bos_token bool             = false
llama_model_loader: - kv  12:                       tokenizer.ggml.model str              = gpt2
llama_model_loader: - kv  13:                      tokenizer.ggml.tokens arr[str,51200]   = ["!", "\"", "#", "$", "%", "&", "'", ...
llama_model_loader: - kv  14:                  tokenizer.ggml.token_type arr[i32,51200]   = [1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, ...
llama_model_loader: - kv  15:                      tokenizer.ggml.merges arr[str,50000]   = ["Ġ t", "Ġ a", "h e", "i n", "r e",...
llama_model_loader: - kv  16:                tokenizer.ggml.bos_token_id u32              = 50296
llama_model_loader: - kv  17:                tokenizer.ggml.eos_token_id u32              = 50297
llama_model_loader: - kv  18:               general.quantization_version u32              = 2
llama_model_loader: - type  f32:  259 tensors
llama_model_loader: - type q8_0:  194 tensors
llm_load_vocab: mismatch in special tokens definition ( 910/51200 vs 944/51200 ).
llm_load_print_meta: format           = GGUF V3 (latest)
llm_load_print_meta: arch             = phi2
llm_load_print_meta: vocab type       = BPE
llm_load_print_meta: n_vocab          = 51200
llm_load_print_meta: n_merges         = 50000
llm_load_print_meta: n_ctx_train      = 2048
llm_load_print_meta: n_embd           = 2560
llm_load_print_meta: n_head           = 32
llm_load_print_meta: n_head_kv        = 32
llm_load_print_meta: n_layer          = 32
llm_load_print_meta: n_rot            = 32
llm_load_print_meta: n_embd_head_k    = 80
llm_load_print_meta: n_embd_head_v    = 80
llm_load_print_meta: n_gqa            = 1
llm_load_print_meta: n_embd_k_gqa     = 2560
llm_load_print_meta: n_embd_v_gqa     = 2560
llm_load_print_meta: f_norm_eps       = 1.0e-05
llm_load_print_meta: f_norm_rms_eps   = 0.0e+00
llm_load_print_meta: f_clamp_kqv      = 0.0e+00
llm_load_print_meta: f_max_alibi_bias = 0.0e+00
llm_load_print_meta: n_ff             = 10240
llm_load_print_meta: n_expert         = 0
llm_load_print_meta: n_expert_used    = 0
llm_load_print_meta: rope scaling     = linear
llm_load_print_meta: freq_base_train  = 10000.0
llm_load_print_meta: freq_scale_train = 1
llm_load_print_meta: n_yarn_orig_ctx  = 2048
llm_load_print_meta: rope_finetuned   = unknown
llm_load_print_meta: model type       = 3B
llm_load_print_meta: model ftype      = Q8_0
llm_load_print_meta: model params     = 2.78 B
llm_load_print_meta: model size       = 2.75 GiB (8.51 BPW)
llm_load_print_meta: general.name     = Phi2
llm_load_print_meta: BOS token        = 50296 '<|im_start|>'
llm_load_print_meta: EOS token        = 50297 '<|im_end|>'
llm_load_print_meta: LF token         = 128 'Ä'
llm_load_tensors: ggml ctx size =    0.35 MiB
llm_load_tensors: offloading 32 repeating layers to GPU
llm_load_tensors: offloading non-repeating layers to GPU
llm_load_tensors: offloaded 33/33 layers to GPU
llm_load_tensors:      ROCm0 buffer size =  2686.46 MiB
llm_load_tensors:        CPU buffer size =   132.81 MiB
............................................................................................
llama_new_context_with_model: n_ctx      = 2048
llama_new_context_with_model: freq_base  = 10000.0
llama_new_context_with_model: freq_scale = 1
llama_kv_cache_init:      ROCm0 KV buffer size =   640.00 MiB
llama_new_context_with_model: KV self size  =  640.00 MiB, K (f16):  320.00 MiB, V (f16):  320.00 MiB
llama_new_context_with_model: graph splits (measure): 3
llama_new_context_with_model:      ROCm0 compute buffer size =   147.00 MiB
llama_new_context_with_model:  ROCm_Host compute buffer size =     9.00 MiB
AVX = 1 | AVX_VNNI = 0 | AVX2 = 1 | AVX512 = 0 | AVX512_VBMI = 0 | AVX512_VNNI = 0 | FMA = 1 | NEON = 0 | ARM_FMA = 0 | F16C = 1 | FP16_VA = 0 | WASM_SIMD = 0 | BLAS = 1 | SSE3 = 1 | SSSE3 = 1 | VSX = 0 |
08:11:56-083810 INFO     LOADER: llama.cpp
08:11:56-084582 INFO     TRUNCATION LENGTH: 2048
08:11:56-085278 INFO     INSTRUCTION TEMPLATE: Alpaca
08:11:56-085904 INFO     Loaded the model in 4.26 seconds.
08:11:56-086636 INFO     Loading the extension "openai"
08:11:56-153090 INFO     OpenAI-compatible API URL:

                         http://0.0.0.0:5000

08:11:56-154806 INFO     Loading the extension "gallery"
Running on local URL:  http://0.0.0.0:7860

To create a public link, set `share=True` in `launch()`.
```

Looks good! If you get permission issues check your permissions on the /opt/text-generation-webui folder

```
sudo systemctl daemon-reload
sudo systemctl enable textgen.service
sudo systemctl start textgen.service
sudo systemctl status textgen.service
● textgen.service
     Loaded: loaded (/etc/systemd/system/textgen.service; enabled; vendor preset: enabled)
     Active: active (running) since Mon 2024-01-22 07:53:18 UTC; 5s ago
   Main PID: 2499 (bash)
      Tasks: 8 (limit: 18905)
     Memory: 254.5M
        CPU: 4.054s
     CGroup: /system.slice/textgen.service
             ├─2499 /bin/bash /opt/text-generation-webui/start_linux.sh
             ├─2512 python one_click.py
             ├─2517 /bin/sh -c ". \"/opt/text-generation-webui/installer_files/conda/etc/profile.d/conda.sh\" && conda activate \"/opt/text-generation-webui/installer_files/env\" && python server.py  --listen --api --model Home-3B-v2.q8_0.gguf --n-gpu-layers 33"
             └─2520 python server.py --listen --api --model Home-3B-v2.q8_0.gguf --n-gpu-layers 33

systemd[1]: Started textgen.service.
```

Add the GPU server IP and port (5000) to your Home-LLM integration.

Select text-generation-webui API in the dropdown then hit SUBMIT.

![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/185ed546-e4ec-4c14-9a16-e77a9b9532c3)

The GPU Server IP that is running text-generation-webui

API Hostname*: 10.0.0.42

The backend port for text-generation-webui 5000 is the default for text-generation-webui (not 7860 which is the webui) 

API Port*: 5000 

Model name has to be the exact same as it looks in the dropdown model list in text-generation-webui

Model Name*: Home-3B-v2_q8_0.gguf

Chat completions endpoint is how the URLs get formed to post data to text-generation-webui without this checked it will use /v1/completions which is depreciated for /v1/chat/completions (Both currently still work but let's think about the future).

[X] Use chat completions endpoint

API Key doesn't matter can type anything or nothing at all.

API key: na

![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/09839a50-4f1c-4cc8-8593-9aa241dfdfe2)

Select SUBMIT and hope for Success!

![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/09726308-5ecc-400c-86c6-404bdac5a1e4)

![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/a0d0ca6a-b2c4-4014-95db-ba30a45c2e27)



## Wyoming-Whisper Container for Docker ##

[![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/9da6de53-17b8-4cf9-98a0-fd20e3651a9f)](https://my.home-assistant.io/redirect/supervisor_addon?addon=core_whisper)

#### Manual Docker Config ####
```yaml
# docker run -it -p 10300:10300 -v /path/to/local/data:/data rhasspy/wyoming-whisper --model tiny-int8 --language en
# Whisper Container for Docker
version: "3.9"
services:
  wyoming-whisper:
    image: rhasspy/wyoming-whisper:latest
    labels:
      - "com.centurylinklabs.watchtower.enable=true"
      - "com.centurylinklabs.watchtower.monitor-only=false"
    container_name: whisper
    volumes:
      - /etc/localtime:/etc/localtime:ro
      - /var/lib/libvirt/images/usb001/docker-storage/wyoming-whisper/data:/data
      #- /etc/asound.conf:/etc/asound.conf
    mac_address: "de:ad:be:ef:5e:34"
    networks:
      - dhcp
    devices:
      - /dev/snd:/dev/snd
    ports:
      - 10300:10300 # http
    command: --model tiny-int8 --language en
    restart: "no"
networks:
  dhcp:
    name: "dbrv100"
```

## Wyoming-Piper Docker Container ##

[![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/9da6de53-17b8-4cf9-98a0-fd20e3651a9f)](https://my.home-assistant.io/redirect/supervisor_addon?addon=core_piper)

#### Manual Docker Config ####
```yaml
# docker run -it -p 10200:10200 -v /path/to/local/data:/data rhasspy/wyoming-piper --voice en-us-lessac-low
# Piper Container for Docker
version: "3.9"
services:
  wyoming-piper:
    image: rhasspy/wyoming-piper:latest
    labels:
      - "com.centurylinklabs.watchtower.enable=true"
      - "com.centurylinklabs.watchtower.monitor-only=false"
    container_name: piper
    volumes:
      - /etc/localtime:/etc/localtime:ro
      - /var/lib/libvirt/images/usb001/docker-storage/wyoming-piper/data:/data
    mac_address: "de:ad:be:ef:5e:35"
    networks:
      - dhcp
    ports:
      - 10200:10200 # http
    command: --voice en-us-lessac-low
    restart: "no"
networks:
  dhcp:
    name: "dbrv100"
```

## Wyoming-OpenWakeWord Container for Docker ##

[![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/9da6de53-17b8-4cf9-98a0-fd20e3651a9f)](https://my.home-assistant.io/redirect/supervisor_addon?addon=core_openwakeword)

#### Manual Docker Config ####
```yaml
# Wyoming-OpenWakeWord Container for Docker
version: "3.9"
services:
  wyomingopenwakeword:
    container_name: "wyoming-openwakeword"
    image: "rhasspy/wyoming-openwakeword:latest"
    labels:
      - "com.centurylinklabs.watchtower.enable=true"
      - "com.centurylinklabs.watchtower.monitor-only=false"
    volumes:
      - /etc/localtime:/etc/localtime:ro
      - /var/lib/libvirt/images/usb001/docker-storage/wyoming-openwakeword/data:/data
    mac_address: "de:ad:be:ef:5e:36"
    hostname: "wyomingopenwakeword"
    networks:
      - dhcp
    ports:
      - 10400:10400
    devices:
      - /dev/snd:/dev/snd
    command: --model 'ok_nabu' --model 'hey_jarvis' --model 'hey_rhasspy' --model 'hey_mycroft' --model 'alexa' --preload-model 'ok_nabu'    
    # --uri 'tcp://0.0.0.0:10400'
    restart: "no"
    #restart: unless-stopped
    #privileged: true
    #network_mode: host

networks:
  dhcp:
    name: "dbrv100"
    external: true
```

## Wyoming Protocol ##

[![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/a57099e9-f9f0-4042-839e-2feebbd14580)](https://my.home-assistant.io/redirect/config_flow_start?domain=wyoming)

![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/a37d4221-d1f6-48a4-9218-e281709e72a4)

In the Wyoming protocol integration click "Add Entry" and enter the IP and Port of the three Whisper/Piper/OpenWakeWord dockers so 10200, 10300, 10400, and their corresponding IP.

Then go under the settings -> voice assistant -> add assistant.

Name your assistant.

Select your conversation agent like Home Assistant or your Home LLM, etc.

Home LLM has some options as well that you can play with

![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/b8d7fd9f-bc51-4315-89c0-0c5c587bfb8c)

Then drop down the boxes for Faster-Whisper, Piper, and OpenWakeWord. Select the settings you like for each.

![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/4c606111-b55b-4a4c-9d9c-ff0309b9c93a)

Then give it a test

![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/673b4cea-828c-4df3-9156-7bd3d122385e)

![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/5e92a85b-599c-43b0-93b7-bc29ecd269b1)

## Home Assistant Groups ##

Currently, the model doesn't support turning off areas so you require this integration to get that working with the model

https://www.home-assistant.io/integrations/group/
 
[![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/a57099e9-f9f0-4042-839e-2feebbd14580)](https://my.home-assistant.io/redirect/config_flow_start?domain=group)

![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/50e38f97-a88a-4e91-bfc6-ea9c25e5adcc)

![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/7d32e0da-0f26-482f-b226-2ccbc8a19b13)

## Wire-Pod Docker ##

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

## Compile ##

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
//agentID := "AgentIDHere" // Replace with your agent_id (Can get this with the dev assist console in YAML view or try the name)
```

### AgentID ###

If you want Vector to use a specific Agent ID you setup for him under assist manager. You need to uncomment these three lines and add your agentID := "################################" etc, to the middle one. Otherwise, it will use the default one.

![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/c155cd51-62f9-414b-9b48-b07f6b3e009c)

Grab the agent_id:

![image](https://github.com/NonaSuomy/WirePodVectorCommandHomeAssistant/assets/1906575/06388744-a45b-4fce-aa7a-7e3a852043af)

```C
//AgentID string `json:"agent_id"`

//agentID := "AgentIDHere" // Replace with your agent_id (Can get this with the dev assist console in YAML view or try the name)

//AgentID: agentID,
```

**Note:** _If you delete and remake your agent in the HA setup it will generate a new AgentID which you will then have to recompile into this and restart the server/vector._

### Compiling ###
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
Then turn it back on...
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
Haven't received a conn check from 00###### in 15 seconds, will ping jdocs on next check
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

### Using a different Agent ID (HomeLLM) ###
https://github.com/acon96/home-llm
```
Bot 00###### Stream type: OPUS
(Bot 00######, Vosk) Processing...
Using general recognizer
(Bot 00######) End of speech detected.
Bot 00###### Transcribed text: assist what does the fox say
Bot 00###### matched plugin Home Assistant Control, executing function
Bot 00###### plugin Home Assistant Control, response The fox does not say.
This is a custom intent or plugin!
Bot 00###### request served.
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
