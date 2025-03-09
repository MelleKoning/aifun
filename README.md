# aifun
aifun is to run an llm locally

## ollama

The idea of the `docker-compose.yaml` file is to have a singular way of starting ollama and openwebui.

The volume mapping is to ensure the pulled models are shared whenever we would start ollama from the commandline in linux.

## startup command

```bash
docker-compose up
```

The above command will startup ollama and openwebui and keep the console open. This way, you can see the incoming HTTP requests that openwebui is executing against ollama. 

Especially when using certain tools, like websearch, within openwebui, it is nice to see what requests ollama and openwebui are executing "under the hood". Sometimes requests fail and all of the logging that the docker containers are doing makes it a bit more clear what is happening when `openwebui` is exercising requests with `ollama`, and how the diferent LLM models are being exercised by `ollama`. For example, at startup, you can see if `ollama` was able to find NVIDIA hardware, as if not, it will fall back to CPU which makes performance far worse.

## Trouble shooting

Sometimes starting up the docker-compse does not work because another ollama instance is already running on your system. This is not surprising as ollama does have start commands, but does not have stop commands. Essentially, the ollama tool is a client tool that ensures that "ollama serve" is running, and when it does, it does not stop (ref: version that was used so far does not have a stop command for the server)

Investigate what other instance of the ollama server is running on port 11434 with

`sudo netstat -peanut | grep ollama`

you can kill other ollama processes or try to stop those via several commands

To stop a server that runs with ubuntu or linux mint service:

```bash
systemctl stop ollama
```

to stop an already running LLM:

```bash
ollama stop
```

to stop an ollama docker instance, you have to lookup the docker-id of the ollama instance and stop and or remove it:

```bash
docker ps -a
docker stop <ollama docker instace>
docker rm <ollama docker instance>
```


etc

