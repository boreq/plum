<template>
    <div class="category-traffic"
         :class="{ dimmed: !checked }"
         :title="active ? 'Click to remove this filter.' : 'Click to display only this category.'"
         v-on:click="click">
        <div class="header">
            <div class="title">
                {{ title }}
            </div>
            <i class="icon fas" :class="checked ? 'fa-eye' : 'fa-eye-slash'"></i>
        </div>
        <div class="metrics">
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
    </div>
</template>
<script lang="ts" src="./CategoryTraffic.ts"></script>
<style scoped lang="scss" src="./CategoryTraffic.scss"></style>
