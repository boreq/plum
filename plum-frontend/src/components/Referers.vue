<template>
    <div>
        <Table :header="header" :rows="rows" :clickable="true" v-on:click-row="clickRow($event)"></Table>
    </div>
</template>
<script lang="ts" src="./Referers.ts"></script>
<style scoped lang="scss" src="./Referers.scss"></style>
