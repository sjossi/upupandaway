# UP Modifier

A tool to deeply unpack extracted .up files, a random file format I found
online. Tested only on unpacked .up files.

It's an update format with instructions in an .ini file. The goal is to unpack
it in such a way that the target filesystem is as closely recreated as possible
to allow for better inspection. There are steps that operate on flash memory,
installing bootloaders etc. These need to be analyzed separately, but since the
update process unpacks them to /tmp, it should all make sense.

Still work in progress, as I ran out of time before a deadline and it's on the
backburner again.


# Code

All features are written as library first and the CLI is just convenience.
The CLI will lag behind and potentially never receive all features, since it's
not the primary focus of this work and I'd rather focus my time on analyzing
firmware than polishing the utility. I also expect people that do research to
be able to handle a bit of code, even in a foreign language (I try my best to
write clean and readable code.)


# TODO

*Library*

* [ ] Standardize extraction folder
* [ ] Unpack/copy binary.ini files
* [ ] Functions for special steps (rootfs unpack)
    * Simulate shellscript based on research?
    * Execute in isolated environment (Docker?) and then move files?
* [ ] Check hashes where provided
* [ ] Repacking

*CLI*

* [ ] Unpack all
