<script lang="ts" setup>
import EventTable from "components/EventTable.vue";
</script>

<template>
  <q-page class="flex column no-wrap">
    <div class="q-pa-md non-selectable">
      <q-chip outline color="primary" text-color="white" icon="event">
        Add to calendar
      </q-chip>
      <q-chip outline color="teal" text-color="white" icon="bookmark">
        Bookmark
      </q-chip>
      <q-chip outline color="orange" text-color="white" icon-right="star">
        Star
      </q-chip>
      <q-chip
        outline
        square
        color="red"
        text-color="white"
        icon="alarm"
        label="Set alarm"
      />
      <q-chip
        outline
        square
        color="deep-orange"
        text-color="white"
        icon="directions"
      >
        Get directions
      </q-chip>
    </div>
    <div class="col q-px-md q-pb-md" style="overflow: hidden">
      <EventTable />
    </div>
  </q-page>
</template>

<style>
.q-page {
  height: 100vh;
  min-height: 0;
}
</style>
