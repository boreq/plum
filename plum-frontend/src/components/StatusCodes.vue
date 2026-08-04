<template>
    <div>
        <Table :header="header" :rows="rows" :clickable="true" v-on:click-row="clickRow($event)"></Table>
    </div>
</template>
<script lang="ts" src="./StatusCodes.ts"></script>
<style scoped lang="scss" src="./StatusCodes.scss"></style>
