# PixBlur

PixBlur is a web app game where the user tries to guess the image within 30 seconds.
The image is updated every second with a slightly less blurry version.
A new image is released daily (currently not an implemented feature)

## Notes
This is meant to be a pet project to learn Go. While I've used my good friend Claude to assist with the MVP (primarily css layout). I intend to code everything myself (with the exception of autocomplete) and only consult AI for implementation ideas/improvements. This project is meant to improve my knowledge and capabilities of Go and web programming.

As of 1/25/2025, this is at best an MVP solution. The image and win state are currently fixed, and there are plenty of improvements to be made.

## TODO
    - ~~Change Gaussian Blur functionalty to generate all images once (currently it generates them for each gameplay)~~
    - Change win state to be fully maintained from server (maybe... mostly on server as is)
    - Provide a score to the user (tbd since it already shows a timer)
    - ~~Add a start game button (game currently start on load)~~
    - Offload image generation to be run from a route (currently generates on server start)
    - Implement basic CMS for managing daily games (adding pic, target word, etc.)