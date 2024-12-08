#!/bin/bash

# Пути к файлам
POSTGRESQL_CONF="/var/lib/postgresql/data/postgresql.conf"  # Укажите путь к postgresql.conf
CUSTOM_CONF="/db_configurations/custom.conf"          # Укажите путь к custom.conf

# Проверяем, существует ли файл custom.conf
#if [ -f "$CUSTOM_CONF" ]; then
#    echo "Файл $CUSTOM_CONF найден. Добавляем его содержимое в $POSTGRESQL_CONF."
#
#    # Добавляем содержимое custom.conf в postgresql.conf
#    echo -e "\n# Включение параметров из custom.conf" >> "$POSTGRESQL_CONF"
#    cat "$CUSTOM_CONF" >> "$POSTGRESQL_CONF"
#
#    echo "Содержимое $CUSTOM_CONF успешно добавлено в $POSTGRESQL_CONF."
#else
#    echo "Ошибка: Файл $CUSTOM_CONF не найден. Проверьте путь и повторите попытку."
#    exit 1
#fi

echo "include = '$CUSTOM_CONF'" >> $POSTGRESQL_CONF