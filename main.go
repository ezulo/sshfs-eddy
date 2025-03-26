package main

import (
    "log"
    "errors"

    "github.com/gotk3/gotk3/glib"
    "github.com/gotk3/gotk3/gtk"
)

const TESTLINE0 = "jimmy_md0  22 rsa /home/eduardo/.ssh/id_ed21559 /mnt/md0  /mnt/sshfs/jimmy/md0 "
const TESTLINE1 = "jimmy_sdb1 22 rsa /home/eduardo/.ssh/id_ed21559 /mnt/sdb1 /mnt/sshfs/jimmy/sdb1"

// GTK: Incrementing enum for column headers
const (
    COLUMN_ID int = iota
    COLUMN_HOSTNAME
    COLUMN_PORT
    COLUMN_AUTH_TYPE
    COLUMN_AUTH_KEY
    COLUMN_REMOTE_DIR
    COLUMN_LOCAL_DIR
    COLUMN_STATE
)

// Enum for mountpoint state
const (
    MOUNT_STATE_UNKNOWN int = iota
    MOUNT_STATE_UNMOUNTED   = 1
    MOUNT_STATE_MOUNTED     = 2
)

type mountpoint struct {
    id          string
    hostname    string
    port        int
    auth_type   string
    auth_key    string
    remote_dir  string
    local_dir   string
    state       int
}

func stateToString(state int) (string, error) {
    switch (state) {
    case MOUNT_STATE_UNKNOWN:
        return "Unknown", nil
    case MOUNT_STATE_UNMOUNTED:
        return "Not Mounted", nil
    case MOUNT_STATE_MOUNTED:
        return "Active", nil
    default:
        return "", errors.New("Invalid state")
    }
}

func createColumn(title string, id int) *gtk.TreeViewColumn {
    cellRenderer, err := gtk.CellRendererTextNew()
    if err != nil {
        log.Fatal("Unable to create text cell renderer:", err)
    }
    column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
    if err != nil {
        log.Fatal("Unable to create cell column:", err)
    }
    return column
}

func setupTreeView() (*gtk.TreeView, *gtk.ListStore) {
    treeView, err := gtk.TreeViewNew()
    if err != nil {
        log.Fatal("Unable to create tree view:", err)
    }
    treeView.AppendColumn(createColumn("ID",                COLUMN_ID         ))
    treeView.AppendColumn(createColumn("Hostname",          COLUMN_HOSTNAME   ))
    treeView.AppendColumn(createColumn("Port",              COLUMN_PORT       ))
    treeView.AppendColumn(createColumn("Auth Type",         COLUMN_AUTH_TYPE  ))
    treeView.AppendColumn(createColumn("Auth Key",          COLUMN_AUTH_KEY   ))
    treeView.AppendColumn(createColumn("Remote Directory",  COLUMN_REMOTE_DIR ))
    treeView.AppendColumn(createColumn("Local Directory",   COLUMN_LOCAL_DIR  ))
    treeView.AppendColumn(createColumn("State",             COLUMN_STATE      ))

    listStore, err := gtk.ListStoreNew(
        glib.TYPE_STRING, 
        glib.TYPE_STRING, 
        glib.TYPE_INT,
        glib.TYPE_STRING, 
        glib.TYPE_STRING,
        glib.TYPE_STRING, 
        glib.TYPE_STRING,
        glib.TYPE_STRING,
    )
    if err != nil {
        log.Fatal("Unable to create list store:", err)
    }
    treeView.SetModel(listStore)
    return treeView, listStore
}

func addRow(
    listStore *gtk.ListStore,
    mp mountpoint,
) {
    state_string, err := stateToString(mp.state)
    if err != nil {
        log.Fatal("Error converting state:", err)
    }
    iter := listStore.Append()
    err = listStore.Set(iter,
        []int{
            COLUMN_ID,
            COLUMN_HOSTNAME,
            COLUMN_PORT,
            COLUMN_AUTH_TYPE,
            COLUMN_AUTH_KEY,
            COLUMN_REMOTE_DIR,
            COLUMN_LOCAL_DIR,
            COLUMN_STATE,
        },
        []interface{}{
            mp.id, 
            mp.hostname, 
            mp.port, 
            mp.auth_type, 
            mp.auth_key, 
            mp.remote_dir, 
            mp.local_dir, 
            state_string,
        })
    if err != nil {
        log.Fatal("Cannot add row:", err)
    }
}

func setupWindow(title string) *gtk.Window {
    win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    if err != nil {
        log.Fatal("Unable to create window:", err)
    }
    win.SetTitle(title)
    win.Connect("destroy", func() {
        gtk.MainQuit()
    })
    win.SetPosition(gtk.WIN_POS_CENTER)
    win.SetDefaultSize(900, 300)
    return win
}

func getMountPoints() []mountpoint {
    var mp []mountpoint
    mp = append(mp,
        mountpoint {
            id: "jimmy_mediapool",
            hostname: "jimmy",
            port: 22,
            auth_type: "rsa",
            auth_key: "/path/to/key",
            remote_dir: "/mediapool",
            local_dir: "/mnt/jimmy_mediapool",
            state: MOUNT_STATE_UNKNOWN,
        },
    )
    mp = append(mp,
        mountpoint {
            id: "jimmy_sdb1",
            hostname: "jimmy",
            port: 22,
            auth_type: "rsa",
            auth_key: "/path/to/key",
            remote_dir: "/mnt/sdb1",
            local_dir: "/mnt/jimmy_sdb1",
            state: MOUNT_STATE_UNKNOWN,
        },
    )
    return mp
}

func main() {
    // Create Gtk Application, change appID to your application domain name reversed.
    const appID = "org.gtk.example"
    gtk.Init(nil)
    win := setupWindow(appID)
    treeView, listStore := setupTreeView()
    win.Add(treeView)
    mp := getMountPoints()
    for i := range mp {
        addRow(listStore, mp[i])
    }
    win.ShowAll()
    gtk.Main()
}

