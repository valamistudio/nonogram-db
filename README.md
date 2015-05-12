## Nonograms

Nonograms (or griddlers, paint by numbers, picross, hanjie, or a million other names) are picture logic puzzles in which cells in a grid must be colored or left blank according to numbers at the side of the grid to reveal a hidden picture.

It's a popular type of puzzle, and most platforms have several apps that let users solve and/or create nonograms.

## Licenses

However, nonogram puzzles that are freely distributable are actually hard to find.  Many web sites let users create nonograms, but few (none?) allow those nonograms to be placed under any sort of license that allows redistribution.

The few open source nonogram programs that do ship with puzzles either have a paltry amount by default or ship clearly dubious-in-origin puzzles (i.e. ripped straight from a commercial game).

This database is an attempt to improve on the status quo.

## Database

All puzzle files can be found in the `db` directory.  Related puzzles will be collected in subdirectories, along with relevant information in README.md files.

## Format

An additional complication for writing a nonogram app is that there are many many formats ([at least 26](http://webpbn.com/export.cgi)) for puzzle files.  None are particularly bizarre or innovative.  They are all just different.

This database uses one format.  It's an existing format (Steve Simpson's `non` format) extended slightly (to support license information and multi-color puzzles).  Details are in FORMAT.md.
