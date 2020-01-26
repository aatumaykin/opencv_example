# import the necessary packages
import argparse
import cv2  # Import the OpenCV library


def main():
    """
    Main method of the program.
    """

    # construct the argument parser and parse the arguments
    ap = argparse.ArgumentParser()
    ap.add_argument("-i", "--image", required=True, help="Path to the image to be thresholded")
    args = vars(ap.parse_args())

    # load the image and convert it to grayscale
    image = cv2.imread(args["image"])

    cv2.imshow("Image", image)
    cv2.waitKey(0)


if __name__ == '__main__':
    main()
    cv2.destroyAllWindows()
