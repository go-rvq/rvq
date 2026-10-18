<template>
  <video
    class="video"
    ref="video"
    :class="[{ hidden: !deviceId }, ...classList.split(' ')]"
    autoplay
    v-bind="$attrs"
  />
  <canvas ref="canvas" style="display: none" />
  <div ref="shutter" class="shutter"></div>
</template>

<style scoped>
.video {
  height: -moz-available;
  height: -webkit-fill-available;
  height: fill-available;
}

.hidden {
  display: none;
}

.shutter {
  opacity: 0;
  transition: all 30ms ease-in;
  position: fixed;
  height: 0%;
  width: 0%;
  pointer-events: none;

  background-color: black;

  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  -webkit-transform: translate(-50%, -50%);
}

.shutter.on {
  opacity: 1; /* Shutter Transparency */
  height: 100%;
  width: 100%;
}
</style>

<script setup lang="ts">
import {drawRotated, screenRotation} from './utils'
import {onBeforeUnmount, onMounted, ref, watch} from 'vue'


class Local {
  inited:boolean = false
  videoStream?: MediaStream
}

declare const navigator: Navigator

const props = defineProps({
  // if should remember last camera
  rememberDevice: {
    type: Boolean,
    default: true
  },
  // try to use these device instead of first one, if the camera label has any keyword from this list
  preferCamerasWithLabel: {
    type: Array<string>,
    default: ['back', 'usb']
  },
  // class list of video element
  classList: {
    type: String,
    default: ''
  },
  // constraints that will be passed to getUserMedia, you can specify preferred resolution, facing direction etc.
  constraints: {
    type: Object,
    default: () => ({
      video: { width: { ideal: 2560 }, height: { ideal: 1440 }, facingMode: 'environment' }
    })
  },
  // if device has gyroscope and the device is rotated (for example in landscape mode), this will try to rotate the image
  tryToRotateImage: {
    type: Boolean,
    default: true
  },
  // output image
  imageType: {
    type: String,
    default: 'image/jpeg'
  },
  // will store the last used device in the local storage if rememberDevice is enabled
  rememberDeviceTokenName: {
    type: String,
    default: '_vwl_device_id'
  },
  // if should automatically start and select the best device depending to preferCamerasWithLabel and constraints, or selects first device
  autoStart: {
    type: Boolean,
    default: true
  },
  audio: {
    type: Boolean,
    default: false
  },
  shutterEffect: {
    type: Boolean,
    default: true
  }
})

const $emit = defineEmits([
  'clear',
  'stop',
  'start',
  'pause',
  'resume',
  'error',
  'unsupported',
  'init',
  'photoTaken'
]),
  local = new Local()

const deviceId = ref(''),
  cameras = ref<InputDeviceInfo[]>([]),
  video = ref<HTMLVideoElement | null>(null),
  canvas = ref<HTMLCanvasElement | null>(null),
  shutter = ref<HTMLDivElement | null>(null)

function init() {
  if (!local.inited) {
    if (!deviceId.value && props.autoStart) {
      start()
    }

    $emit('init', deviceId.value)
    local.inited = true
  }
}

function start() {
  if (deviceId.value) {
    loadCamera(deviceId.value)
  } else {
    // check if there is any remembered device and if so, use them
    let rememberedDevice: string | null = null
    try {
      rememberedDevice = window.localStorage.getItem(props.rememberDeviceTokenName)
    } catch {
      /* no storage (a private window) */
    }
    if (
      rememberedDevice &&
      props.rememberDevice &&
      cameras.value.find((el) => el.deviceId === rememberedDevice)
    ) {
      deviceId.value = rememberedDevice
    } else if (cameras.value.length > 1) {
      for (const label of props.preferCamerasWithLabel) {
        const camera = cameras.value.find((el) => el.label.toLowerCase().indexOf(label) !== -1)
        if (camera) {
          deviceId.value = camera.deviceId
          break
        }
      }
    }
    // nothing found, use first if there is any
    if (!deviceId.value && cameras.value.length > 0) {
      deviceId.value = cameras.value[0].deviceId
    }
  }
  $emit('start')
}

function loadCameras(): Promise<InputDeviceInfo[]> {
  // resolved even when the devices cannot be listed (unsupported emitted)
  return new Promise((resolve) => {
    navigator.mediaDevices
      .enumerateDevices()
      .then((deviceInfos:any) => {
        for (let i = 0; i !== deviceInfos.length; ++i) {
          let deviceInfo = deviceInfos[i]
          // need to include only devices with proper deviceId (as without permission the deviceId is = '')
          if (
            deviceInfo.deviceId &&
            deviceInfo.kind === 'videoinput' &&
            cameras.value.find((el:any) => el.deviceId === deviceInfo.deviceId) === undefined
          ) {
            cameras.value.push(deviceInfo)
          }
        }
      })
      .then(() => {
        if (!local.inited && cameras.value.length > 0) {
          if (!deviceId.value && props.autoStart) {
            start()
          }
          $emit('init', deviceId.value)
          local.inited = true
        }
        resolve(cameras.value)
      })
      .catch((error:any) => {
        $emit('unsupported', error)
        resolve(cameras.value)
      })
  })
}

function changeCamera(devId: string) {
  if (deviceId.value !== devId) {
    deviceId.value = devId
    return // will be recalled due to watcher
  }
  stop()
  if (devId) {
    loadCamera(devId)
  }
}

function loadCamera(deviceId: string) {
  navigator.mediaDevices
    .getUserMedia(buildConstraints(deviceId))
    .then((stream:any) => {
      (video.value as any).srcObject = stream

      if (props.rememberDevice) {
        try {
          window.localStorage.setItem(props.rememberDeviceTokenName, deviceId)
        } catch {
          /* no storage (a private window) */
        }
      }
    })
    .catch((err: any) => $emit('error', err))
}

function buildConstraints(deviceId: string=''): MediaStreamConstraints {
  // a copy: the device chosen is never written into the props' constraints
  // (a later call would ask that device, gone or not)
  const c: any = { audio: false, ...props.constraints }
  const video = typeof c.video === 'object' && c.video !== null ? { ...c.video } : {}
  if (deviceId) {
    video.deviceId = { exact: deviceId }
    // a device chosen: its facing is its own
    delete video.facingMode
  }
  c.video = Object.keys(video).length ? video : true
  return c as MediaStreamConstraints
}

function testMediaAccess() {
  navigator.mediaDevices
    .getUserMedia(buildConstraints())
    .then((stream: any) => {
      let tracks = stream.getTracks()
      tracks.forEach((track: any) => {
        track.stop()
      })
      loadCameras()
    })
    .catch((err: any) => $emit('error', err))
}

function legacyGetUserMediaSupport() {
  return (constraints?: MediaStreamConstraints) => {
    let getUserMedia = ((navigator as any).getUserMedia ||
      (navigator as any).webkitGetUserMedia ||
      (navigator as any).mozGetUserMedia ||
      (navigator as any).msGetUserMedia ||
      (navigator as any).oGetUserMedia) as (constraints?: MediaStreamConstraints, resolve?:(value:MediaStream)=>void, reject?:(reason:any)=>void)=>Promise<MediaStream>

    if (!getUserMedia) {
      return Promise.reject(new Error('getUserMedia is not implemented in this browser'))
    }
    return new Promise<MediaStream>((resolve, reject) => {
      getUserMedia.call(navigator, constraints, resolve, reject)
    })
  }
}

function setupMedia() {
  if (navigator.mediaDevices === undefined) {
    (navigator as any).mediaDevices = {}
  }
  if (navigator.mediaDevices.getUserMedia === undefined) {
    navigator.mediaDevices.getUserMedia = legacyGetUserMediaSupport()
  }
  testMediaAccess()
}

function clear() {
  const v = video.value as HTMLVideoElement | null
  const stream = v?.srcObject as MediaStream | null | undefined
  stream?.getTracks().forEach((track) => track.stop())
  if (v) v.srcObject = null
  $emit('clear')
}

function stop() {
  if (video.value?.srcObject) {
    clear()
  }
  $emit('stop')
}
function pause() {
  if ((video.value as any).srcObject) {
    video.value?.pause()
  }
  $emit('pause')
}

function resume() {
  if ((video.value as any).srcObject) {
   video.value?.play()
  }
  $emit('resume')
}

async function  takePhoto() {
  let v = video.value as HTMLVideoElement
  let cv = canvas.value as HTMLCanvasElement
  cv.height = v.videoHeight
  cv.width = v.videoWidth
  let ctx = cv.getContext('2d') as CanvasRenderingContext2D

  drawRotated(
    v,
    cv,
    ctx,
    props.tryToRotateImage ? screenRotation() : 0
  )

  let image_data_url = cv.toDataURL(props.imageType)
  cv.toBlob((blob) => {
    if (props.shutterEffect) {
      shutter.value?.classList.add('on')
      setTimeout(
        () => {
          shutter.value?.classList.remove('on')
        },
        30 * 2 + 45
      )
    }
    $emit('photoTaken', { blob, image_data_url })
  }, props.imageType)
}

onBeforeUnmount(()=> {
  stop()
})

onMounted(()=>{
  setupMedia()
})

watch(deviceId, (deviceId)=> {
  changeCamera(deviceId)
})

defineExpose({
  loadCameras,
  changeCamera,
  loadCamera,
  start,
  stop,
  pause,
  resume,
  init,
  takePhoto,
  cameras
})
</script>
