import subprocess, os, sys, shutil, json


def report_compile_pdf():
    latex_command = 'pdflatex'
    tex_file_path = 'a.tex'
    compile_command = [latex_command, '-shell-escape', tex_file_path]

    process = subprocess.Popen(compile_command, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    stdout, stderr = process.communicate()

    if process.returncode == 0:
        print("Файл успешно скомпилирован")
    else:
        print(f"Ошибка при компиляции файла:\n{stderr.decode('utf-8')}")


def get_code(path, extensions=['.go']):
    result = ''
    for root, dirs, files in os.walk(path):
        for file_name in files:
            if file_name.endswith(tuple(extensions)):
                file_path = os.path.join(root, file_name)
                with open(file_path, 'r', encoding='utf-8') as file:
                    content = file.read()
                    content = split_content(content, 125)
                    result += f'Файл {file_name}\n\\begin{{minted}}{{go}}\n{content}\\end{{minted}}\n\n'
    return result


def split_content(content, max_length=100):
    content_to_split = content
    lines = [content_to_split[i:i + max_length] for i in range(0, len(content_to_split), max_length)]
    return '\n'.join(lines)


def report_make(path, data):
    # Получение текущей директории
    latex_template = '../latex-template/template.tex'

    try:
        with open(latex_template, 'r', encoding='utf-8') as file:
            content = file.read()
    except FileNotFoundError:
        print(f"Файл '{latex_template}' не найден.")
        return
    except Exception as e:
        print(f"Ошибка при чтении файла: {e}")
        return

    code = get_code(path)
    content = content.replace('{{.labNumber}}', path[path.find('/'):].replace('/', '').replace('lab', ''))
    content = content.replace('{{.code}}', code)
    content = content.replace('{{.labName}}', data['labName'])
    content = content.replace('{{.labVariant}}', data['labVariant'])
    content = content.replace('{{.statement}}', data['statement'])


    # записываем файл, который будем компилировать
    with open('a.tex', 'w', encoding='utf-8') as new_file:
        new_file.write(content)


def clean_after_compile(file_prefix='a'):
    file_extensions = ['.aux', '.log', '.out', '.tex']
    minted_folder = f'_minted-{file_prefix}'

    for file_ext in file_extensions:
        file_name = f'{file_prefix}{file_ext}'
        if os.path.exists(file_name):
            os.remove(file_name)

    os.remove('emblem.png')

    if os.path.exists(minted_folder):
        for root, dirs, files in os.walk(minted_folder, topdown=False):
            for file in files:
                os.remove(os.path.join(root, file))
            for dir in dirs:
                os.rmdir(os.path.join(root, dir))
        os.rmdir(minted_folder)


# скопируем в нашу директорию нужные нам файлы
def prepare_for_compile(path):
    # Копирование emblem.png
    shutil.copy('../latex-template/emblem.png', os.getcwd())

    file_path = os.path.join(path, 'config.json')
    try:
        with open(file_path, 'r', encoding='utf-8') as file:
            config_data = json.load(file)
    except FileNotFoundError:
        print(f"Файл '{file_path}' не найден.")
    except json.JSONDecodeError as e:
        print(f"Ошибка при чтении JSON файла: {e}")
    except Exception as e:
        print(f"Произошла ошибка: {e}")

    files_to_copy = ['bla.png', 'blabla.png']
    for file_name in files_to_copy:
        source_file = os.path.join(path, file_name)
        if os.path.exists(source_file):
            shutil.copy(source_file, os.getcwd())
            print(f"Файл {file_name} скопирован.")
        else:
            print(f"Файл {file_name} не найден по указанному пути {source_file}.")

    return config_data


def compress_folder_to_zip(path):
    base_name = path[path.find('/') + 1:].replace('/', '') # Получаем имя папки для создания имени архива
    shutil.make_archive(base_name, 'zip', path)  # Создаем архив

    # Перемещаем архив в текущую директорию
    archive_name = f"{base_name}.zip"
    shutil.move(archive_name, os.path.join(os.getcwd(), archive_name))
    print(f"Папка '{path}' сжата в архив '{archive_name}' в текущей директории.")



if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Использование: python3 main.py <путь_к_директории>")
        sys.exit(1)

    directory_path = sys.argv[1]

    data = prepare_for_compile(directory_path)
    report_make(directory_path, data)
    report_compile_pdf()
    compress_folder_to_zip(directory_path)

    clean_after_compile()

