import tkinter as tk
from tkinter import filedialog, messagebox
import base64
import os

# Variables to store file and folder paths
file_path = None
output_folder = None

# Function to select an input file
def select_file():
    global file_path
    file_path = filedialog.askopenfilename()
    if file_path:
        file_label.config(text=f"Selected File: {os.path.basename(file_path)}")

# Function to select the output folder
def select_output_folder():
    global output_folder
    output_folder = filedialog.askdirectory()
    if output_folder:
        folder_label.config(text=f"Output Folder: {output_folder}")

# Function to convert file content to base64
def convert_to_base64():
    if not file_path or not output_folder:
        messagebox.showerror("Error", "Please select both a file and an output folder.")
        return
    
    try:
        # Read the file content as bytes
        with open(file_path, 'rb') as file:
            file_content = file.read()

        # Convert to base64 string
        base64_string = base64.b64encode(file_content).decode('utf-8')
 ### REWRITE!!!
        # Prepare the new file name with .ote extension
        file_name = os.path.basename(file_path)
        new_file_name = os.path.splitext(file_name)[0] + ".ote"
        new_file_path = os.path.join(output_folder, new_file_name)

        # Write the base64 string to the new file
        with open(new_file_path, 'w') as new_file:
            new_file.write(base64_string)

        # Display success message
        messagebox.showinfo("Success", f"File converted and saved as:\n{new_file_path}")
    except Exception as e:
        # Handle errors
        messagebox.showerror("Error", f"An error occurred: {str(e)}")

# Function to create the GUI
def create_gui():
    global file_label, folder_label

    # Create the main window
    root = tk.Tk()
    root.title("OpenTunnelV2 File Converter")
    root.geometry("600x400")  # Increased window size

    # Create a label for "OpenTunnelV2" at the top-center
    label_title = tk.Label(root, text="OpenTunnelV2", font=("Helvetica", 18), anchor="center")
    label_title.pack(pady=1)
    label_subtitle = tk.Label(root, text="Host", font=("Helvetica", 12), anchor="center")
    label_subtitle.pack(pady=1)

    # Create a button to select the input file
    select_button = tk.Button(root, text="Select File", command=select_file)
    select_button.pack(pady=10)

    # Label to show the selected file
    file_label = tk.Label(root, text="Selected File: None", anchor="w")
    file_label.pack(pady=5)


    # Create a button to select the output folder
    folder_button = tk.Button(root, text="Select Folder", command=select_output_folder)
    folder_button.pack(pady=5)
    
    # Create a label for output folder selection
    folder_label = tk.Label(root, text="Output Folder: None", anchor="w")
    folder_label.pack(pady=5)
    
    # Create a "Go" button to trigger the conversion
    go_button = tk.Button(root, text="Go", command=convert_to_base64)
    go_button.pack(pady=20)

    # Start the Tkinter event loop
    root.mainloop()

# Main execution
if __name__ == "__main__":
    create_gui()
