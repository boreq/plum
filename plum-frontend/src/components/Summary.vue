<template>
    <div class="summary">
        <div class="metric">
            <div class="label">
                Visits
            </div>
            <div class="value">
                {{ visits }}
            </div>
        </div>
        <div class="metric">
            <div class="label">
                Hits
            </div>
            <div class="value">
                {{ hits }}
            </div>
        </div>
        <div class="metric">
            <div class="label">
                Bytes sent
            </div>
            <div class="value">
                {{ bytesSent }}
            </div>
        </div>
    </div>
</template>
<script lang="ts" src="./Summary.ts"></script>
<style scoped lang="scss" src="./Summary.scss"></style>
