import sys
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service
from selenium.common.exceptions import NoSuchElementException
from webdriver_manager.chrome import ChromeDriverManager
import time

def get_kb_update_info(kb_number):
    # Set up the WebDriver
    service = Service(ChromeDriverManager().install())
    options = webdriver.ChromeOptions()
    options.add_argument('--headless')
    options.add_argument('--disable-gpu')
    driver = webdriver.Chrome(service=service, options=options)

    # Navigate to the Microsoft Update Catalog
    url = f"https://www.catalog.update.microsoft.com/Search.aspx?q={kb_number}"
    print(f"Navigating to URL: {url}")
    driver.get(url)

    # Wait for the page to load
    time.sleep(5)

    # Print the page source for debugging
    page_source = driver.page_source
    print(f"Page source loaded. Length: {len(page_source)} characters")

    description = ""
    size = ""
    try:
        table = driver.find_element(By.ID, "ctl00_catalogBody_updateMatches")
        rows = table.find_elements(By.TAG_NAME, "tr")[1:]  # Skip the header row
        print(f"Found {len(rows)} rows in the results table")
        
        if rows:
            columns = rows[0].find_elements(By.TAG_NAME, "td")
            title_element = columns[1].find_element(By.TAG_NAME, "a")
            size = columns[6].text.strip()
            title = title_element.text.strip()

            # Click the title link to open the pop-up
            print(f"Clicking on title: {title}")
            title_element.click()
            time.sleep(10)  # Wait for the pop-up to load

            # Switch to the pop-up window
            driver.switch_to.window(driver.window_handles[-1])

            # Scrape the pop-up content
            description_element = driver.find_element(By.ID, "ScopedViewHandler_desc")
            description = description_element.text.strip()
            print(f"Scraped description: {description}")

            # Close the pop-up window and switch back to the main window
            driver.close()
            driver.switch_to.window(driver.window_handles[0])
            time.sleep(2)  # Wait for the main window to be focused

    except NoSuchElementException:
        print("No results found for the given KB number.")
    except Exception as e:
        print(f"Error processing rows: {e}")

    driver.quit()
    return description, size

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python script.py <KB_NUMBER>")
        sys.exit(1)

    kb_number = sys.argv[1]
    description, size = get_kb_update_info(kb_number)
    if description and size:
        print(f"Description: {description}\nSize: {size}")
    else:
        print("No description or size found or failed to retrieve data.")
