package application;

import java.util.concurrent.CopyOnWriteArrayList;

import javafx.scene.Scene;
import javafx.scene.shape.Circle;

public class CircleFactory implements Runnable {

	private final int nCircle;
	private final Scene scene;
	private final CopyOnWriteArrayList<Circle> circles;

	public CircleFactory(int nCircle, CopyOnWriteArrayList<Circle> circles,
			Scene scene) {
		this.nCircle = nCircle;
		this.scene = scene;
		this.circles = circles;
	}

	@Override
	public void run() {
		
		for (int i = 0; i < nCircle; i++) {
			// circle holder
			Circle c;
			// Condition for which the circle is valid
			boolean condition = true;

			// if the circle overlaps to other circles then recreate it
			do {
				// generate random circle
				c = Main.circleFactory(scene.getWidth(), scene.getHeight());

				// set the condition to true again
				condition = true;

				// check against all previous circles if there are previous
				// circles
				if (circles.isEmpty()) {
					break;
				}

				// check all circle if is overlapped
				for (final Circle circle : circles) {
					if (Main.overlap(c, circle)) {
						condition = false;
						break;
					}
				}
				// keep on generating the next circle until it doesn't overlap
				// with any of the previous ones
			} while (condition == false);

			// at this point the circle is valid so it's added to the list
			circles.add(c);
		}
	}

}
