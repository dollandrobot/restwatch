<script setup>
// import EventTable from "components/EventTable.vue";
import { GetMessages } from "app/wailsjs/go/main/App.js";
import { ref } from "vue";

const columns = [
  {
    name: "name",
    required: true,
    label: "Event",
    align: "left",
    field: (row) => row.rawMessage,
    format: (val) => `${val}`,
    sortable: true,
  },
  {
    name: "receivedAt",
    required: true,
    label: "Received At",
    align: "left",
    field: (row) => row.receivedAt,
    format: (val) => `${val}`,
    sortable: true,
  },
];

const messages = ref([]);

const onReceiveMessage = (message) => {
  console.log("Received message:", message);
  messages.value.push(message);
};

const loadMessages = () => {
  GetMessages()
    .then((m) => {
      messages.value = m;
    })
    .catch((error) => {
      console.error("Error loading messages:", error);
    });
};

window.runtime.EventsOn("messageReceived", onReceiveMessage);
</script>

<template>
  <q-page class="flex flex-center">
    <q-btn color="primary" label="Refresh" @click="loadMessages" />
    <!--    <EventTable />-->
    <q-table title="Event" :rows="messages" :columns="columns" row-key="name" />
  </q-page>
</template>
