# It's Pong!
The old school Atari Pong game remade in terminal with Golang.  

# General Info
For best experience, make your terminal fullscreen or zoom out.  
The game currently works **only on Linux systems** due to platform-dependent key detection.  
For a similar reason, this game **must be run with superuser privileges**.  

To interact with terminal, the [tcell](https://github.com/gdamore/tcell) package API is used.  
For keyboard events detection, the [keylogger](https://github.com/MarinX/keylogger) package API is used.  

## Why sudo?
It is a compromise that had to be made. Since you cannot detect keyboard press and release key events through terminal, package [keylogger](https://github.com/MarinX/keylogger) is used to do it. Basically, it reads input from `/dev/input/event*` which requires root access. It is a drawback, but without it, the game would not be playable.

# Compile & Run
To compile this project into an executable, run:
```shell
go get github.com/spinzed/tpong
```
You should be able to run the game via command `sudo tpong`.  
**Note:** make sure than `$GOPATH/bin` is in your path, otherwise you won't be able to run the command from anywhere else but from `$GOPATH/bin`.  

Alternative is to git clone this repo, cd into it and run:
```shell
sudo go run .
```

# Game Screenshot
![Image Not Found](https://i.ibb.co/rvmNys6/game.jpg)

# Controls
|  Key  |     Action    |
|:-----:|:-------------:|
|   W   |   Player1 Up  |
|   S   |  Player1 Down |
|   ↑   |   Player2 Up  |
|   ↓   |  Player2 Down |
| Space |     Start     |
|   Q   |      Quit     |
|   P   |  Toggle Pause |
|   R   | Restart Round |

# Further Plans & Todos
Eventually, this project may be expanded:
1. that it can be hosted via HTTP API
2. that it is compatible with online play

# Keep in Mind
This project is primarily an exercise for me to learn more about the Go language. That being said, **any constructive feedback is appreciated**!
