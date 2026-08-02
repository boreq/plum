<template>
    <div id="#app" class="just-site">
        <div class="just-content">
            <div class="just-container-wrapper">
                <div class="just-container">
                    <router-view/>
                </div>
            </div>
        </div>
        <footer class="just-footer">
            <div class="just-container-wrapper">
                <div class="just-container">
                    <a href="https://github.com/boreq/plum">Github</a>
                </div>
            </div>
        </footer>
    </div>
</template>

<style lang="scss">
    html {
        font-family: 'Roboto', sans-serif;
        font-size: 12px;
    }

    // Sticky footer
    .just-site-wrapper {
        height: 100%;
        margin: 0;

        .just-site {
            height: 100%;
            margin: 0;

            display: flex;
            flex-direction: column;

            .just-content {
                flex: 1 0 auto;
            }

            .just-footer {
                flex-shrink: 0;
            }
        }
    }

    .just-container-wrapper {
        padding: 1em;

        .just-container {
            max-width: 1100px;
            margin-left: auto;
            margin-right: auto;
        }
    }

    // Other
    html, body {
        background-color: $background-color;
    }

    footer {
        color: $box-dimmed-text-color;

        a {
            color: inherit;
            text-decoration: none;
        }

        a:hover {
            color: $primary-color;
        }
    }
</style>
