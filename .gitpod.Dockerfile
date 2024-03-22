FROM gitpod/workspace-go:2024-03-20-07-19-19

# Install dependencies (ngrok)
RUN wget https://bin.equinox.io/c/bNyj1mQVY4c/ngrok-v3-stable-linux-amd64.tgz --no-check-certificate && \
    tar -zxvf ngrok-v3-stable-linux-amd64.tgz && \
    sudo mv ngrok /usr/bin/ngrok && \
    sudo chmod 755 /usr/bin/ngrok && \ 
    # air
    go install github.com/cosmtrek/air@latest
