# chkb [![Go Report Card](https://goreportcard.com/badge/github.com/MetalBlueberry/chkb)](https://goreportcard.com/report/github.com/MetalBlueberry/chkb) [![Coverage Status](https://coveralls.io/repos/github/MetalBlueberry/chkb/badge.svg?branch=master)](https://coveralls.io/github/MetalBlueberry/chkb?branch=master)

chkb turns a regular keyboard intro a fully programmable keyboard.

So you basically get a **ch**eap programmable **k**e**y**board.

It has been inspired by [QMK](https://docs.qmk.fm/#/) firmware and [kmonad](https://github.com/david-janssen/kmonad).

## Usage

This applies to the current preview version, this will change in the future but I will try to keep this info updated.

> ATTENTION: if you manage to block your keyboard, hold ESC and the application will quit. This is a safe guard and cannot be disabled and is handled before any key event.

create a `keys.yaml` with at least a base layer keyMap.

```yaml
layers:
  base:
    keyMap:

```

find the input file from your keyboard. This files are located at `/dev/input/by-id/`. The id should be enough to tell which one is your keyboard. if not, you can `sudo cat /dev/input/by-id/file` and see if you can see data when you type on your keyboard.

To avoid sudo, you can add your user to the input group.

```sh
sudo usermod -a -G input $USER
```

chkb needs access to `/dev/uinput` in order to generate key events. You can create a rule for this. so you don't need sudo.

> Extracted from here: https://github.com/bendahl/uinput

```sh
echo KERNEL==\"uinput\", GROUP=\"$USER\", MODE:=\"0660\" | sudo tee /etc/udev/rules.d/99-$USER.rules
sudo udevadm trigger
```

Run `chkb` from the directory containing this file providing as first parameter the path to your device. If everything works, you should see a message `push layer base` which means that everything is running. You can start typing to see some captured events.

```sh
./chkb /dev/input/by-id/usb-Logitech_USB_Receiver-if01-event-kbd
# TIP add -v flag for debug info
```

## Remap keys

This example remaps the CAPSLOCK to LEFTMETA. This means that pressing caps will behave as if you've pressed leftmeta.

```yaml
layers:
  base:
    keyMap:
      KEY_CAPSLOCK:
        Map:
          - keyCode: KEY_LEFTMETA
```

## Push/Pop layers

A layer modifies how your keys behave. Base layer is pushed at startup and cannot be removed. You can put maps on this layer to push other layers so you can extend the functionality. The layers are push into a stack and readed from top to down. Once the event is handled by a layer, the bottom layers do not receive it.

### swap AB keys

This example shows how to create a layer that temporally swaps keys AB

```yaml
layers:
  base:
    keyMap:
      KEY_CAPSLOCK:
        # Tap means to press and release in less than 200ms
        Tap:
          # PushLayer requires the layer name to push that must match with yaml key
          - action: PushLayer
            layerName: swapAB
  # definition of swapAB layer
  swapAB:
    keyMap:
      # swap A
      KEY_A:
          - keyCode: KEY_B
      # swap B
      KEY_B:
          - keyCode: KEY_A
      # remove the layer if caps is tapped again
      KEY_CAPSLOCK:
        Tap:
          - action: PopLayer
            layerName: swapAB
```

## Tap / Hold

You can capture special events like tap/hold and perform custom actions

### CAPSLOCK on SHIFT hold

This will tap CAPSLOCK if you hold LEFTSHIFT

```yaml
layers:
  base:
    keyMap:
      KEY_LEFTSHIFT:
        Hold:
          - action: Tap
            keyCode: KEY_CAPSLOCK
```

## Multiple maps

There are cases where you want to do multiple actions with a single input event. You may have noticed that the Keyboard events are a list, This means you can use as many as you need. 

### Pop layer and push another

This example shows how to enable a control layer that allows you to jump to another layer

```yaml
layers:
  base:
    keyMap:
      KEY_CAPSLOCK:
        Tap:
          - action: PushLayer
            layerName: control
        Map:
          - keyCode: KEY_LEFTMETA

  ## Intermediate layer
  control:
    keyMap:
      KEY_CAPSLOCK:
        ## Go back to base
        Tap:
          - action: PopLayer
            layerName: control
        ## Ensure key still works as a meta key
        Map:
          - keyCode: KEY_LEFTMETA
      KEY_A:
        # If tap A, pop this layer and push arrows layer
        Tap:
          - action: PopLayer
            layerName: control
          - action: PushLayer
            layerName: arrows
        # Block a key so it doesn't type anything
        Map:
          - action: Nil
  arrows:
    keyMap:
      KEY_CAPSLOCK:
        ## Go back to base by poping arrows layer
        Tap:
          - action: PopLayer
            layerName: arrows
        ## Ensure key still works as a meta key
        Map:
          - keyCode: KEY_LEFTMETA

      ## Put arrows on hjkl vim style
      KEY_J:
        Map:
          - keyCode: "KEY_DOWN"
      KEY_K:
        Map:
          - keyCode: "KEY_UP"
      KEY_H:
        Map:
          - keyCode: "KEY_LEFT"
      KEY_L:
        Map:
          - keyCode: "KEY_RIGHT"

```

## Available key events

key events are used to capture events from your keyboard and map them into other keyboard events. 

```yaml
layers:
  base:
    keyMap:
      KEY_LEFTSHIFT:
        Hold: <--- This is the key event
          - action: Tap
            layerName: KEY_CAPSLOCK
```

- Map: Forwards the up/down events to a different key
- Down: On key down
- Up: On key up
- Tap: press and release in less than 200ms, no other keys are pressed in between
- DoubleTap: tap twice!
- Hold: keep the key down for more than 200ms or press another key while it is down.

## Available keyboard events

These events are the ones generated by the keyMaps.

```yaml
layers:
  base:
    keyMap:
      KEY_LEFTSHIFT:
        Hold: 
          - action: Tap  <--- This is a Keyboard Event
            layerName: KEY_CAPSLOCK 
```

- Nil: does nothing
- Map: Forwards the up/down events to a different key // Default behaviour
  - keyCode: Key used
- Down: Simulate key up
  - keyCode: Key used
- Up: Simulate key down
  - keyCode: Key used
- Tap: press and release in less than 200ms
  - keyCode: Key used
- DoubleTap: tap twice!
  - keyCode: Key used
- Hold: press until the keyboard repeats the event
  - keyCode: Key used
- PushLayer: push a layer into the stack
  - layerName: name of the layer
- PopLayer: removes a layer from the stack
  - layerName: name of the layer


## Keycodes

If you tap a key, you will see the code printed in the chkb log. This way you can find the name of any key from your keyboard.

To find extra keys, you can read ./pkg/chkb/ecodes.go and use anything starting with `KEY_`. This keys are just standard keycodes and you can also find them in linux source code /usr/include/linux/input-event-codes.h

## TODO

- [x] integrate with cobra cli
- [ ] implement testing run mode. just to see keypresses and final maps
- [ ] allow to specify the location of the config file
- [ ] easy autoshift functionality for layers
- [ ] examples
- [ ] Allow to disable layerFile
- [ ] Implement RESETKEY to start-stop the app in the background