# hblock
hyper block command line tools and SDK library. 

Main page: http://www.hyperblock.org

_Currently support VM format 'qcow2', and 'lvm' will be supported in the near future_

## Installation

1. Install command __'qcow2-img'__ and create a soft link.   
compile srouce code from https://github.com/hyperblock/qemu/tree/hyper-block, use 'make qcow2-img' to build.

2. Build hblock from the source.  
    *  pre-install need
        * go-flags      github.com/jessevdk/go-flags
        * yaml          gopkg.in/yaml.v2  
        * go.uuid       github.com/satori/go.uuid
        * libguestfs    libguestfs.org/guestfs
        
    * use __'go build hb.go'__ and create a soft link.

3. Need a __WebDAV__ server if use some remote options.

## How to use 

    hb <command> [options]

	Note. use 'hb <command> -h' to see detail.

    ======== support commands =======

	init            create empty backingfile
	config          get and set global options
	clone           clone a repository from remote or local path
	remote          manage set of tracked repositories
	rebase          reapply volume's backingfile and parent-layer 
	branch          list,create or delete branches
	checkout        switch branches or restore volume
	commit          record volume's changes
	reset           reset current HEAD to the specified state
	pull            fetch from and integrate with another repository of a local branch
	push            update remote repository
    list            list backingfiles in current workspace
	show            show backingfile's detail
    log             show commit logs
	
* ### __init__  
    Usage:
      hb init <template name> [OPTIONS]
    
    Application Options:  
          --size=   [required] Disk size(M/G) of template.
                    eg.
                    hblock init template0 --size=500M -f qcow2.
    
      -o=           [optional] output volume name.    
      -f, --format= 'qcow2' of 'lvm'.

* ### __config__  
    Usage:
      hb config [OPTIONS]  
    Application Options:
          --global= [user.name|user.email] set global configuration.
          --get=    <name>    Get value : <name>

* ### __clone__
    Usage:
      hb clone <repo path> [OPTIONS]
    
    Application Options:
      -l, --layer=       Checkout <layer> instead of the HEAD    
          --hardlink   use local hardlinks.    
      -b, --branch=      Clone the specified <branch> instead of default ('master').    
      -n, --no-checkout  No checkout of HEAD is performed after clone is complete.

* ### __remote__  
    Usage:
      hb remote <volume> [OPTIONS]
    
    Application Options:
      -v, --verbose
      -a, --add      <name> <url>    Add a new remote-host to local remote-host list.
      -d, --remove   <name>	Delete a host from local remote-host list.
          --rename   <old_name> <new_name>	 Rename an exsiting host name.
          --set-url

* ### __rebase__
    Usage:
      hb rebase <volume_name> [OPTIONS]
    
    Application Options:
      -b, --backingfile= <backingfile>
      -l, --layer=       <layer>

* ### __branch__
    Usage:
      hb branch [OPTIONS]
    
    Application Options:
          --list         list branch names.    
      -a, --all          list both remote-tracking and local branches.    
      -m, --move=        <exist_branch> <new_branch> move/rename a branch.    
      -t, --backingfile= required if use '-m'    
      -v, --volume=    

* ### __checkout__
    Usage:
      hb checkout [OPTIONS]    
    Application Options:
      -v, --vol=         <volume_name> <layer | branch> Specify the volume name which needs to be
                         update(restore).    
      -t, --backingfile= <backingfile> <layer | branch> Create a new volume from <backingfile>.    
      -o, --output=      <output_volume_path>.    
      -b, --branch=      <branch> <volume_name> Create a new branch of specified volume.    
      -f, --force.

* ### __commit__
    Usage:
      hb commit <volume name> [OPTIONS]
    
    Application Options:
      -m=         commit message
          --uuid= set uuid by manual instead of auto-generate.

* ### __reset__
    Usage:
        hb reset <volume> [<commit_uuid>] | [HEAD point]	reset <volume> and discard changes.
    	eg.
    		hb reset volume0 3f2ed7		reset 'volume' to specified commit 3f2ed7
    		hb reset volume0 HEAD^^		reset 'volume' to the last 2 commits
    		hb reset volume0 HEAD~5		reset 'volume' to the last 5 commits    	

* ### __pull__
    Usage:
      hb pull <remote> <branch> [OPTIONS]
        
    Application Options:
      -v, --volume=

* ### __push__
    Usage:
      hb push <repository> <refspec> [OPTIONS]
    
    Application Options:
      -v, --volume= <volume>

* ### __log__
    Usage:
      hb log <volume name> [OPTIONS]
    
* ### __list__
    Usage:
        hb list
    	hb list <dir_path>
    
* ### __show__
    Usage:
        hb show <backing file> 	show backing file details.




