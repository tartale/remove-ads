# RMads

RMads - app for detecting and removing ads from a video stream

## Summary

Given a set of markers that describe time segments where commercials have
been detected within a video stream, RMads will format and execute an
FFMpeg command to remove the commercials from the stream.

The goals of this project are:
- Create a basic CLI app that can import skip data and remove the
indicated sections from the video
- Create a UI to allow interactive preview and correction of the
markers being used to cut the ads
- Create an ML model that can be trained to improve the detection
algorithms

## Credits
