# Easy Open

Open urls in your OS default browser, with configurable shortcuts
# Usage
```
easyopen google.com
or
easyopen yourcmd yourparam1 ...
```

# Install
Download the latest package for your OS:
```
wget https://github.com/Easy-Infra-Ltd/easy-open/releases/download/v0.0.1.2/easyopen-v0.0.1.2-linux-amd64.tar.gz
```
Extract the files locally:
```
tar -xzf easyopen-v0.0.1.2-linux-amd64.tar.gz
```
Run the install script to get setup:
```
sudo ./install
```

# Configure
If you followed the install steps above inside your `~/.config` folder you will have a file named `easyopen.cmds.json` that should look something like this:
```
[
    {
        "name": "yt",
        "url": "youtu.be/:1"
    },
    {
        "name": "ytsearch",
        "url": "youtube.com/results?search_query=:1"
    },
    {
        "name": "google",
        "url": "google.com/search?q=:1"
    },
    {
        "name": "so",
        "url": "stackoverflow.com/search?q=:1"
    }
]
```

This allows you to setup shortcuts and pass parameters, they will be numbered in order of entry after the command and can be re-used through the url multiple times. For example:
```
example.com/:1/:1-:2
```
