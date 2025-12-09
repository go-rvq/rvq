<template>
  <v-card height="100%" style="display: flex; flex-direction: column">
    <v-toolbar density="compact" v-if="closable" style="justify-content: center">
      <IconButton
        :title="translate.loadCameras"
        density="compact"
        size="x-large"
        icon="mdi-camera-retake"
        @click="loadCameras"
      />
      <div>
        <v-select
          hide-details
          density="compact"
          :items="cameras"
          item-title="label"
          item-value="deviceId"
          v-model="deviceId"
          @update:model-value="setCamera"
        ></v-select>
      </div>
      <IconButton
        :title="translate.flipCamera"
        density="compact"
        size="x-large"
        icon="mdi-camera-flip"
        @click="flipCamera"
        v-if="cameras.length > 1"
      />
      <template v-slot:prepend>
        <v-btn disabled />
      </template>
      <template v-slot:append>
        <IconButton :title="translate.close" icon="mdi-close" @click="closeHandler" />
      </template>
    </v-toolbar>
    <v-divider />
    <v-card-text
      style="
        overflow: hidden;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 0;
      "
      class="webcam"
    >
      <WebCam
        ref="webcam"
        @init="webcamInit"
        @clear="clear"
        @stop="stop"
        @start="start"
        @pause="pause"
        @resume="resume"
        @error="error"
        @unsupported="unsupported"
        @photoTaken="photoTakenEvent"
        :shutterEffect="fullscreen"
      />
    </v-card-text>
    <v-divider />
    <v-card-actions style="justify-content: center">
      <template v-if="picture.url">
        <IconButton
          :title="translate.clear"
          variant="flat"
          density="compact"
          size="x-large"
          icon="mdi-arrow-left"
          @click="takenCancel"
        />
        <IconButton
          :title="translate.save"
          color="primary"
          variant="flat"
          density="compact"
          size="x-large"
          icon="mdi-check-bold"
          @click="doSave"
        />
      </template>
      <IconButton
        v-else
        :title="translate.takePicture"
        color="primary"
        variant="flat"
        density="compact"
        size="x-large"
        icon="mdi-image"
        @click="takePhoto"
      />
    </v-card-actions>
  </v-card>
</template>

<script lang="ts" setup>
import { onBeforeUnmount, onMounted, useTemplateRef, watch, ref, computed, useAttrs } from 'vue'
import { Translator } from './translator'
import IconButton from '../IconButton.vue'
import { Picture } from './types.js'
import WebCam from './WebCam.vue'

class interval {
  value:unknown = 0
}

const props = defineProps({
    fullscreenState: {
      type: Boolean,
      default: false
    },
    closable: {
      type: Boolean,
      default: false
    },
    translate: {
      type: Translator,
      default: () => new Translator()
    },
    save: {
      type: Function,
      default: null
    }
  }),
  webcam = useTemplateRef('webcam'),
  cameras = computed(():InputDeviceInfo[] => {
    return (webcam.value?.cameras || []) as InputDeviceInfo[]
  }),
  deviceId = ref(''),
  fullscreen = ref(false),
  photoTaken = ref(false),
  photoFailed = ref(false),
  reloadCamInterval = new interval(),
  picture = ref(new Picture()),
  $emit = defineEmits([
    'clear',
    'stop',
    'start',
    'photoSave',
    'pause',
    'resume',
    'error',
    'unsupported',
    'init',
    'photoTaken',
    'fullscreen',
    'close'
  ]),
  takePhoto = async () => {
    try {
      await webcam.value?.takePhoto()
      photoTaken.value = true
      setTimeout(() => {
        photoTaken.value = false
      }, 500)
    } catch (err) {
      photoFailed.value = true
      setTimeout(() => {
        photoFailed.value = false
      }, 500)
    }
  },
  loadCameras = () => {
    return webcam.value?.loadCameras() as Promise<any[]>
  },
  webcamInit = (devId:string) => {
    deviceId.value = devId
    $emit('init', devId)
  },
  setCamera = () => {
    webcam.value?.changeCamera(deviceId.value)
  },
  flipCamera = () => {
    loadCameras().then(() => {
      // flipping camera will select the next one from the list, but on most device there will be only 2, if < 2 it will not be shown
      if (cameras.value.length > 1) {
        let currentIndex = cameras.value.findIndex((el) => el.deviceId === deviceId.value)
        let newIndex = currentIndex + 1
        if (newIndex >= cameras.value.length) {
          newIndex = 0
        }

        deviceId.value = cameras.value[newIndex].deviceId
        webcam.value?.changeCamera(cameras.value[newIndex].deviceId)
      }
    })
  },
  exit = () => {
    webcam.value?.stop()
  },
  closeHandler = () => {
    exit()
    $emit('close')
  },
  // emits
  clear = () => {
    $emit('clear')
  },
  stop = () => {
    $emit('stop')
  },
  start = () => {
    $emit('start')
  },
  pause = () => {
    $emit('pause')
  },
  resume = () => {
    $emit('resume')
  },
  error = (err:any) => {
    $emit('error', err)
  },
  unsupported = (err:any) => {
    $emit('unsupported', err)
  },
  photoTakenEvent = (arg:any) => {
    picture.value.set(arg.blob as Blob, arg.image_data_url as string)
    webcam.value?.pause()
  },
  doSave = () => {
    $emit('photoTaken', picture.value)
  },
  takenCancel = () => {
    picture.value.clear()
    webcam.value?.resume()
  }

onMounted(() => {
  if (cameras.value.length === 0) {
    // if no camera found, we will try to refresh cameras list each second until there is some camera
    reloadCamInterval.value = setTimeout(() => {
      loadCameras().then((value) => {
        if (value.length > 0) {
          clearInterval(reloadCamInterval.value as number)
          // most likely due to permission, so we init afterwards
          webcam.value?.init()
        }
      })
    }, 10000) as unknown
  }
})
onBeforeUnmount(() => {
  if (reloadCamInterval.value) {
    clearInterval(reloadCamInterval.value as number)
  }
  exit()
})
</script>
