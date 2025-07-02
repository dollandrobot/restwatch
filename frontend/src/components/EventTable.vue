<script lang="ts">
import {
  ref,
  watch,
  onMounted,
  onBeforeUnmount,
  computed,
  nextTick,
} from "vue";
import type { main } from "../../wailsjs/go/models";
import { GetMessages } from "../../wailsjs/go/main/App";
import { useQuasar } from "quasar";

export default {
  setup(props, { emit }) {
    const $q = useQuasar();
    const virtualListRef = ref(null);
    const messages = ref([]);
    const isDark = computed(() => $q.dark.isActive);
    const automaticScrolling = ref(false);
    const scrollToLatest = ref(true);

    const fetchMessages = async () => {
      try {
        const result = await GetMessages();
        messages.value = result;
        executeScroll();
      } catch (error) {
        console.log("error getting existing messages", error);
      }
    };

    const executeScroll = () => {
      if (scrollToLatest.value) {
        automaticScrolling.value = true;
        virtualListRef.value.scrollTo(messages.value.length - 1, "start-force");
        // Delay resetting automaticScrolling to ensure the scroll animation completes
        setTimeout(() => {
          automaticScrolling.value = false;
        }, 300);
      }
    };

    const onVirtualScroll = () => {
      if (!automaticScrolling.value) {
        scrollToLatest.value = false;
      }
    };

    const onReceiveMessage = (message: main.SimpleMessage) => {
      messages.value.push(message);
      if (scrollToLatest.value) {
        void nextTick(() => {
          executeScroll();
        });
      }
    };

    const onRowClick = (row) => {
      scrollToLatest.value = false;
      emit("row-click", row.id);
    };

    onMounted(async () => {
      await fetchMessages();
      virtualListRef.value.scrollTo(0);
      window.runtime.EventsOn("messageReceived", onReceiveMessage);
    });

    onBeforeUnmount(() => {
      window.runtime.EventsOff("messageReceived");
    });

    watch(scrollToLatest, () => {
      if (scrollToLatest.value) {
        executeScroll();
      }
    });

    return {
      messages,
      onRowClick,
      isDark,
      virtualListRef,
      onVirtualScroll,
      executeScroll,
      scrollToLatest,
      settings: ref(false),

      slideVol: ref(39),
      slideAlarm: ref(56),
      slideVibration: ref(63),
    };
  },
};
</script>

<template>
  <div class="row">
    <div class="col"></div>
    <div class="col text-right q-gutter-sm">
      <q-checkbox v-model="scrollToLatest" label="Scroll to latest" />
      <q-icon
        name="settings"
        size="2em"
        class="cursor-pointer"
        @click="settings = true"
      />
    </div>
  </div>
  <q-virtual-scroll
    ref="virtualListRef"
    flat
    bordered
    type="table"
    :items="messages"
    style="max-height: 100%; height: 100%"
    :virtual-scroll-item-size="48"
    :virtual-scroll-sticky-size-start="48"
    @virtual-scroll="onVirtualScroll"
  >
    <template v-slot:before>
      <thead>
        <tr>
          <th
            :class="{ 'th-dark': isDark, 'th-light': !isDark }"
            class="text-left no-pointer-events"
          >
            #
          </th>
          <th
            :class="{ 'th-dark': isDark, 'th-light': !isDark }"
            class="text-left no-pointer-events"
          >
            Content
          </th>
          <th
            :class="{ 'th-dark': isDark, 'th-light': !isDark }"
            class="text-left no-pointer-events"
          >
            Received At
          </th>
        </tr>
      </thead>
    </template>
    <template v-slot:default="{ item: row, index }">
      <tr :key="index" @click="onRowClick(row)">
        <td>#{{ index + 1 }}</td>
        <td>{{ row.content }}</td>
        <td>{{ row.receivedAt }}</td>
      </tr>
    </template>
  </q-virtual-scroll>

  <q-dialog v-model="settings">
    <q-card style="width: 300px" class="q-px-sm q-pb-md">
      <q-card-section>
        <div class="text-h6">Volumes</div>
      </q-card-section>

      <q-item-label header>Media volume</q-item-label>
      <q-item dense>
        <q-item-section avatar>
          <q-icon name="volume_up" />
        </q-item-section>
        <q-item-section>
          <q-slider color="teal" v-model="slideVol" :step="0" />
        </q-item-section>
      </q-item>

      <q-item-label header>Alarm volume</q-item-label>
      <q-item dense>
        <q-item-section avatar>
          <q-icon name="alarm" />
        </q-item-section>
        <q-item-section>
          <q-slider color="teal" v-model="slideAlarm" :step="0" />
        </q-item-section>
      </q-item>

      <q-item-label header>Ring volume</q-item-label>
      <q-item dense>
        <q-item-section avatar>
          <q-icon name="vibration" />
        </q-item-section>
        <q-item-section>
          <q-slider color="teal" v-model="slideVibration" :step="0" />
        </q-item-section>
      </q-item>
    </q-card>
  </q-dialog>
</template>

<style lang="sass" scoped>
.q-virtual-scroll
  overflow-y: auto

// Common header styles mixin
%th-common
  position: sticky
  top: 0
  z-index: 1
  padding: 8px 16px
  transition: all 0.3s ease
  font-weight: 500

.th-light
  @extend %th-common
  background-color: $deep-purple-1
  color: $deep-purple-10

.th-dark
  @extend %th-common
  background-color: $deep-purple-10
  color: $deep-purple-1

td
  padding: 8px 16px

tr
  cursor: pointer

  &:hover
    background-color: rgba(0, 0, 0, 0.05)
    transition: background-color 0.2s ease

// Dark mode specific hover
:deep(.body--dark) tr
  &:hover
    background-color: rgba(255, 255, 255, 0.1)
    transition: background-color 0.2s ease
</style>
