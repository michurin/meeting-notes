#!/bin/sh

jq '[.categories[]|select(.categories)|{(.title):[.categories[].title]}]|add' ozon >ozon.json
