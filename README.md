# Andy's Code Clipper

## About

This is a little tool I created for my own personal use, but if you like it then please feel free to fork and help improve.

This tool is a little note-taking, code clipper thing.  I wanted to be able to make notes in Markdown, with an emphasis on being able to add code snippets, and have the Markdown rendered and the code shown with syntax highlighting.

I also wanted to have a play with [Wails](https://wails.io/), so this was the perfect opportunity!

## How it works

When you start typing a note it will save automatically.  Each note is an individual json file that's stored in a folder called `code-clipper` inside your home directory.

Code tags will be automatically take from fenced code blocks, so if you use:

~~~
```bash
echo "foo"
```
~~~

It will record the tag `bash`.

There's no limit the number of different fences you use in the same note - they'll all be turned into tags.

If you want a title to the note, it'll automatically use the first header 1 when using the format:

~~~
# This will be the title
~~~

Whenever a note is saved, its update time, title, and tags are reflected in the list.  To filter on the tags, just select the tag you want to show.  To view a note, just click on the note's title.

When you select a note to view, it starts off in rendered mode.  To edit, click the `edit` button.

To delete a note, select the one you want to remove and click the `delete` button... But use caution as there is no coming back from that - there is no undo or confirmation!


## Future plans

At some point I'd like to do the following:

- Improve the code \
Right now it's a bit on the hacky side, so I'd hope to improve upon it
- Allow for different themes to be used for the syntax highlighting \
Right now you get Monaki Dark, and that's it!  But ideally I'd allow selecting more
- MMake deleting a note a bit safer \
There is no confirmation for delete, and once you click that delete button, the note is gone for good!  So I want to explore if it's possible to use the recycle bin when deleting
- Explore more of the Wails configuration \
I've only built and tested this on Windows 10, but would like to better improve the experience, try it on a Mac, and so on

