FROM tars.cppbase

RUN apt update; apt install -y wget unzip

RUN wget https://github.com/nvm-sh/nvm/archive/v0.35.1.zip;unzip v0.35.1.zip; cp -rf nvm-0.35.1 $HOME/.nvm

RUN echo 'NVM_DIR="$HOME/.nvm"; [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"; [ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion";' >> $HOME/.bashrc \
    && . $HOME/.bashrc \
    && nvm install v12.13.0

# RUN apt install -y node
COPY /root /
RUN chmod a+x /bin/entrypoint.sh
# COPY /tars-node-agent /usr/local/app/tars/