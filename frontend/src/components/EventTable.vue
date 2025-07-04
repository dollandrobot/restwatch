<script setup lang="ts">
import {
  ref,
  watch,
  onMounted,
  onBeforeUnmount,
  computed,
  nextTick,
} from "vue";
import type { ComponentPublicInstance } from "vue";
import type { QVirtualScroll } from "quasar";
import type { main } from "../../wailsjs/go/models";
import { GetMessages, SaveUserOptions } from "../../wailsjs/go/main/App";
import { useQuasar } from "quasar";

const props = defineProps<{
  userOptions: main.UserOptions;
}>();

const emit = defineEmits(["row-click"]);

const $q = useQuasar();
const virtualListRef = ref<ComponentPublicInstance<typeof QVirtualScroll>>();
const messages = ref<main.SimpleMessage[]>([]);
const isDark = computed(() => $q.dark.isActive);
const automaticScrolling = ref(false);
const scrollToLatest = ref(true);
const showSettings = ref(false);
const settings = ref(props.userOptions);

const fetchMessages = async () => {
  try {
    setAndTrimMessages(await GetMessages());
    executeScroll();
  } catch (error) {
    console.log("error getting existing messages", error);
  }
};

const setAndTrimMessages = (incoming: main.SimpleMessage[]) => {
  const excess = incoming.length - settings.value.maxMessagesToKeep;
  messages.value = excess > 0 ? incoming.slice(excess) : incoming;
};

const onReceiveMessage = (message: main.SimpleMessage) => {
  setAndTrimMessages([...messages.value, message]);
  if (scrollToLatest.value) {
    void nextTick(() => {
      executeScroll();
    });
  }
};

const executeScroll = () => {
  if (scrollToLatest.value && virtualListRef.value) {
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

const onRowClick = (row: main.SimpleMessage) => {
  scrollToLatest.value = false;
  emit("row-click", row.id);
};

const onSaveOptionsClick = async () => {
  await SaveUserOptions(settings.value);
  showSettings.value = false;
};

onMounted(async () => {
  await fetchMessages();
  if (virtualListRef.value) {
    virtualListRef.value.scrollTo(0);
  }
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
</script>

<template>
  <div class="row">
    <div class="col"></div>
    <div class="col text-right q-gutter-sm">
      <q-checkbox v-model="scrollToLatest" label="Scroll to latest" />
      <q-icon
        name="refresh"
        size="2em"
        class="cursor-pointer"
        @click="fetchMessages()"
      />
      <q-icon
        name="settings"
        size="2em"
        class="cursor-pointer"
        @click="showSettings = true"
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

  <q-dialog v-model="showSettings" backdrop-filter="blur(4px) saturate(150%)">
    <q-card style="width: 300px" class="q-px-sm q-pb-md">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">Settings</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <q-card-section class="q-pt-none">
        <q-input
          v-model.number="settings.maxMessagesToKeep"
          type="number"
          filled
          label="Max # of messages to keep"
        />
      </q-card-section>

      <q-card-section class="q-pt-none">
        <q-input
          v-model="settings.defaultEndpoint"
          filled
          label="Incoming messages path"
        />
      </q-card-section>
      <q-card-section class="q-pt-none">
        <q-input v-model="settings.port" filled label="Port" />
      </q-card-section>
      <q-separator />
      <q-card-actions align="right">
        <q-btn label="Ok" flat color="primary" @click="onSaveOptionsClick()" />
      </q-card-actions>
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
