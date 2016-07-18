package ar.com.bestprice.buyitnow;

import android.view.MotionEvent;
import android.view.View;

/**
 * Created by elagiglia on 18/7/16.
 */
public class SwipeDetector implements View.OnTouchListener{

    public static enum Action {
        LR,
        RL,
        TB,
        BT,
        None
    }
    private Action mSwipeDetected = Action.None;
    private float downX, downY, upX, upY;
    private static final int HORIZONTAL_MIN_DISTANCE = 100;
    private static final int VERTICAL_MIN_DISTANCE = 100;

    public boolean swipeDetected(){return mSwipeDetected == Action.None;}


    @Override
    public boolean onTouch(View view, MotionEvent event) {
        switch (event.getAction()){
            case MotionEvent.ACTION_DOWN:
                downX = event.getX();
                downY = event.getY();
                mSwipeDetected = Action.None;
                return false;

            case MotionEvent.ACTION_BUTTON_PRESS:
                break;

            case MotionEvent.ACTION_BUTTON_RELEASE:
                break;

            case MotionEvent.ACTION_MOVE:

                downX = event.getX();
                downY = event.getY();

                float deltaX = downX - upX;
                float deltaY = downY - upY;

                if (Math.abs(deltaX) > HORIZONTAL_MIN_DISTANCE){
                    if(deltaX < 0){
                        mSwipeDetected = Action.LR;
                        return true;
                    } else if (deltaX > 0) {
                        mSwipeDetected = Action.RL;
                        return true;
                    }
                } else {
                    if (Math.abs(deltaY) > VERTICAL_MIN_DISTANCE){
                        if(deltaY < 0){
                            mSwipeDetected = Action.TB;
                            return false;
                        } else if (deltaY > 0) {
                            mSwipeDetected = Action.BT;
                            return false;
                        }
                    }
                }
                return true;
        }

        return false;
    }

}
