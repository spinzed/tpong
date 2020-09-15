# It's Pong!
The old school Atari Pong game remade in terminal with Golang.  

# General Info
For the best experience, make your terminal fullscreen or zoom it out.  

The game currently works **only on Linux systems** due to platform-dependent key detection.  
For a similar reason, this game **must be run with superuser privileges**.  

To interact with terminal, the [tcell](https://github.com/gdamore/tcell) package API is used.  
For keyboard events detection, the [keylogger](https://github.com/MarinX/keylogger) package API is used.  
For command line argument parsing, the [go-flags](https://github.com/jessevdk/go-flags) package API is used.  

## Why sudo?
It is a compromise that had to be made. Since you cannot detect keyboard press and release key events through terminal, package [keylogger](https://github.com/MarinX/keylogger) is used to do it. Basically, it reads input from `/dev/input/event*` which requires root access. It is a drawback, but without it, the game would not be playable.

# Features
- Two player local multiplayer!
- Command line arguments!
- Start menu!
- Round pausing!
- AI!
- Themes!
- ... and more to come!

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

# Game Screenshots
![Image Not Found](https://i.ibb.co/YN1ZJh7/1.jpg)  
![Image Not Found](https://i.ibb.co/zQpjkDw/2.jpg)  
![Image Not Found](https://i.ibb.co/drdpLBX/3.jpg)  

# Controls
### In menus:
|  Key  |         Action         |
|:-----:|:----------------------:|
| Enter | Choose Selected Action |
|   ↑   |    Scroll Actions Up   |
|   ↓   |   Scroll Actions Down  |

### In game:
|  Key  |       Action      |
|:-----:|:-----------------:|
|   W   |     Player1 Up    |
|   S   |    Player1 Down   |
|   ↑   |     Player2 Up    |
|   ↓   |    Player2 Down   |
| Space |       Start       |
|   Q   |        Quit       |
|   P   |    Toggle Pause   |
|   R   |   Restart Round   |
|   A   |     Toggle AI     |
|   T   |    Switch Theme   |
|   B   | Toggle Background |

## Command line arguments:
|  Flag  |                  Description                 | Default Value | Optional |
|:------:|:--------------------------------------------:|:-------------:|:--------:|
| --nobg | Dictates whether background should be hidden |     false     |    yes   |


# Further Plans & Todos
Eventually, this project may be expanded:
1. that it can be hosted via HTTP API
2. that it is compatible with online play  

**Any form of contribution or feedback is appreciated!**
