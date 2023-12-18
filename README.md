photoFS
=======

## Project Goal

photoFS is a tool written in Golang to manage tags with a GUI, save them to an embedded database and create a tag-based folder-structure/filesystem for your photo collection.

This makes it possible to tag a growing photo collection with only one tool and utilize the tags in multiple ways like for example at local folders, NFS/Samba shares, Web-tools like Nextcloud, TV via Kodi Media-centers etc.

It is written in Golang, to get compiled binaries with least possible dependencies to be save for the future and make it also easier usable for non-tech people.

## Project State

This project is in development and is still missing features which are blocking a productional use.

## Planned Features

- Embedded KV database
	- Projects which could be a good fit:
		- [BadgerDB](https://github.com/dgraph-io/badger)
		- [bolt/bbolt](https://github.com/etcd-io/bbolt)
		- [BuntDB](https://github.com/tidwall/buntdb)
		- [bitcask](https://git.mills.io/prologic/bitcask)
		- [NutsDB](https://github.com/nutsdb/nutsdb)
- Hierarchical tags
- Automatic EXIF tags
- Live updated folder-structure/filesystem on tag updates
- ...

## Utilized Golang Projects
- [Golang Fyne Module](https://github.com/fyne-io/fyne) (For the GUI of the tagging client)

## Background Story of the photoFS Project

As I am a passionate photographer I was already using several photo tagging tools in the past.

I saw myself already facing multiple issues with these tools, for example:
- A tool is might not being maintained in the future anymore
- You put a lot of time into tagging, but usually you can utilize these tags only in the same tool you created the tags with.
	- But I would like to tag my photos only once and be able to use these tags then in multiple ways to view my photo collection:
		- Local folders
		- Remote folders (NFS/Samba shares)
		- Custom Websites
		- File/Photo Web-tools like Nextcloud
		- TV via Kodi Media-center
- You could be switching your operating system, which is not compatible anymore
- Many, if not most of these GUI tools do not offer hierarchical tags with multiple dimensions, which could also be viewed in merged views without duplicated tags

For that reason I wrote already back in 2013 photoFS as a GUI tool to create tags, save them to a dedicated database, which has then been used to create a filesystem.

I wrote it in Perl, utilized Tk for the GUI, MongoDB for the database and Fuse for the filesystem.

Now when looking back, the problem with this approach was that it still had too much dependencies to be used for a lifetime solution and by non-tech people as well.

A solution with a compiled binary, would be way better and nowadays I would choose Golang and an embedded database for that use-case.

That is how this Golang re-write project has been born.

## Legal

### Licence
photoFS is open-sourced software licensed under the GNU GPLv3.

### Licence Notice
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

### Copyright Notice
Copyright 2023, Sven Burkard
